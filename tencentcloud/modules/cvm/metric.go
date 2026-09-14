package cvm

import (
	"context"
	"fmt"
	"time"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/monitor/v20180724"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Shared builder + hydrate for all CVM metric tables (namespace QCE/CVM).
func cvmMetricTable(name, description, metricName, valueDescription string, periodSeconds uint64, lookback time.Duration) *plugin.Table {
	granularity := "hourly"
	if periodSeconds == 86400 {
		granularity = "daily"
	}
	return &plugin.Table{
		Name:              name,
		Description:       description,
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "instance_id", Require: plugin.Optional},
				{Name: "timestamp", Require: plugin.Optional, Operators: []string{">=", ">", "<=", "<", "="}},
			},
			Hydrate: func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
				return listCvmMetric(ctx, d, metricName, periodSeconds, lookback)
			},
			Tags: map[string]string{"service": "monitor", "action": "GetMonitorData"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "instance_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the CVM instance.", Transform: transform.FromField("InstanceId")},
			{Name: "metric_name", Type: proto.ColumnType_STRING, Description: fmt.Sprintf("The name of the metric (%s).", metricName), Transform: transform.FromField("MetricName")},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "The metric namespace of the service (QCE/CVM).", Transform: transform.FromField("Namespace")},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp of the data point in ISO 8601 format (UTC).", Transform: transform.FromField("Timestamp").NullIfZero()},
			{Name: "value", Type: proto.ColumnType_DOUBLE, Description: valueDescription, Transform: transform.FromField("Value")},
			{Name: "period", Type: proto.ColumnType_INT, Description: fmt.Sprintf("The statistical period of the data point in seconds (%d for %s).", periodSeconds, granularity), Transform: transform.FromField("Period")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the instance resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("InstanceId")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmMetric(ctx context.Context, d *plugin.QueryData, metricName string, periodSeconds uint64, lookback time.Duration) (any, error) {
	plugin.Logger(ctx).Debug("listCvmMetric", "region", d.EqualsQualString(utils.MatrixKeyRegion), "metric_name", metricName, "quals", d.EqualsQuals)

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	cvmClient := &cvm.Client{}
	if err := utils.InitClient(ctx, d, cvmClient); err != nil {
		plugin.Logger(ctx).Error("listCvmMetric", "client_init_error", err)
		return nil, err
	}

	monitorClient := &monitor.Client{}
	if err := utils.InitClient(ctx, d, monitorClient); err != nil {
		plugin.Logger(ctx).Error("listCvmMetric", "client_init_error", err)
		return nil, err
	}

	requirePublicIP := metricName == "WanIntraffic" || metricName == "WanOuttraffic"
	instanceIds, err := listCvmInstanceIds(ctx, d, cvmClient, requirePublicIP)
	if err != nil {
		return nil, err
	}
	if len(instanceIds) == 0 {
		return nil, nil
	}

	// Default window: [now-lookback, now]
	endTime := time.Now().UTC()
	startTime := endTime.Add(-lookback)
	periodDuration := time.Duration(periodSeconds) * time.Second
	if quals := d.Quals["timestamp"]; quals != nil {
		for _, q := range quals.Quals {
			ts := q.Value.GetTimestampValue().AsTime().UTC()
			switch q.Operator {
			case "=", ">=", ">":
				if ts.After(startTime) {
					startTime = ts
				}
				if q.Operator == "=" && ts.Add(periodDuration).Before(endTime) {
					endTime = ts.Add(periodDuration)
				}
			case "<", "<=":
				if ts.Before(endTime) {
					endTime = ts
				}
			}
		}
	}
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("invalid timestamp range: start time %s must be before end time %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	}

	const batchSize = 10
	for i := 0; i < len(instanceIds); i += batchSize {
		end := min(i+batchSize, len(instanceIds))
		batch := instanceIds[i:end]

		d.WaitForListRateLimit(ctx)
		req := monitor.NewGetMonitorDataRequest()
		req.Namespace = common.StringPtr("QCE/CVM")
		req.MetricName = common.StringPtr(metricName)
		req.Period = common.Uint64Ptr(periodSeconds)
		req.StartTime = common.StringPtr(startTime.Format(time.RFC3339))
		req.EndTime = common.StringPtr(endTime.Format(time.RFC3339))
		instances := make([]*monitor.Instance, 0, len(batch))
		for _, id := range batch {
			instances = append(instances, &monitor.Instance{
				Dimensions: []*monitor.Dimension{
					{Name: common.StringPtr("InstanceId"), Value: common.StringPtr(id)},
				},
			})
		}
		req.Instances = instances

		utils.LogRequest(ctx, "listCvmMetric", req)
		resp, err := monitorClient.GetMonitorDataWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmMetric", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			continue
		}

		namespace := utils.PtrString(req.Namespace)
		for _, dp := range resp.Response.DataPoints {
			if dp == nil {
				continue
			}
			instanceId := cvmMetricInstanceIdFromDimensions(dp.Dimensions)
			for idx, ts := range dp.Timestamps {
				if ts == nil || idx >= len(dp.Values) || dp.Values[idx] == nil {
					continue
				}
				row := cvmMetricRow{
					InstanceId: instanceId,
					MetricName: utils.PtrString(resp.Response.MetricName),
					Namespace:  namespace,
					Timestamp:  floatUnixSecondsToTime(ts),
					Value:      utils.PtrFloat64(dp.Values[idx]),
					Period:     utils.PtrUint64(resp.Response.Period),
					Region:     region,
				}
				if row.InstanceId != "" {
					row.Akas = []string{"tencentcloud:cvm:instance:" + row.InstanceId}
				}
				d.StreamListItem(ctx, row)

				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

func cvmMetricInstanceIdFromDimensions(dimensions []*monitor.Dimension) string {
	for _, dim := range dimensions {
		if dim == nil {
			continue
		}
		if utils.PtrString(dim.Name) == "InstanceId" {
			return utils.PtrString(dim.Value)
		}
	}
	return ""
}

func floatUnixSecondsToTime(ts *float64) *time.Time {
	if ts == nil || *ts == 0 {
		return nil
	}
	t := time.Unix(int64(*ts), 0).UTC()
	return &t
}

// Row Type

type cvmMetricRow struct {
	InstanceId string
	MetricName string
	Namespace  string
	Timestamp  *time.Time
	Value      float64
	Period     uint64
	Region     string
	Akas       []string
}
