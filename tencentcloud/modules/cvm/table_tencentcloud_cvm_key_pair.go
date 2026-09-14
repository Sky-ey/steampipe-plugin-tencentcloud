package cvm

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCvmKeyPair() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cvm_key_pair",
		Description:       "Tencent Cloud CVM key pairs (SSH login credentials) that can be bound to Linux instances for secure login.",
		GetMatrixItemFunc: buildCvmRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "key_id", Require: plugin.Optional},
				{Name: "key_name", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
			},
			Hydrate: listCvmKeyPairs,
			Tags:    map[string]string{"service": "cvm", "action": "DescribeKeyPairs"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("key_id"),
			Hydrate:    getCvmKeyPair,
			Tags:       map[string]string{"service": "cvm", "action": "DescribeKeyPairs"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "key_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the key pair (e.g., skey-11112222).", Transform: transform.FromField("KeyId")},
			{Name: "key_name", Type: proto.ColumnType_STRING, Description: "The name of the key pair.", Transform: transform.FromField("KeyName")},
			{Name: "project_id", Type: proto.ColumnType_INT, Description: "The ID of the project the key pair belongs to. 0 means the default project.", Transform: transform.FromField("ProjectId")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the key pair.", Transform: transform.FromField("Description")},
			{Name: "public_key", Type: proto.ColumnType_STRING, Description: "The plain-text public key of the key pair.", Transform: transform.FromField("PublicKey")},
			{Name: "associated_instance_ids", Type: proto.ColumnType_JSON, Description: "The IDs of the instances associated with the key pair.", Transform: transform.FromField("AssociatedInstanceIds")},
			{Name: "created_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the key pair in ISO 8601 format (UTC).", Transform: transform.FromField("CreatedTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the key pair as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the key pair resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("KeyName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCvmKeyPairs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCvmKeyPairs", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCvmKeyPairs", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64 = 0

	keyIds := utils.SingleIdFromQual(d.EqualsQuals, "key_id")
	filters := buildKeyPairFilters(d.EqualsQuals)

	for {
		d.WaitForListRateLimit(ctx)

		req := cvm.NewDescribeKeyPairsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(keyIds) > 0 {
			req.KeyIds = keyIds
		} else {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listCvmKeyPairs", req)
		resp, err := client.DescribeKeyPairsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listCvmKeyPairs", "request_error", err)
			return nil, err
		}

		if resp == nil || resp.Response == nil {
			break
		}
		keyPairs := resp.Response.KeyPairSet
		total := utils.PtrInt64(resp.Response.TotalCount)

		for _, keyPair := range keyPairs {
			if keyPair == nil {
				continue
			}
			row := toCvmKeyPairRow(keyPair)
			row.Region = region
			d.StreamListItem(ctx, row)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(offset, int64(len(keyPairs)), total) {
			break
		}
		offset += int64(len(keyPairs))
	}

	return nil, nil
}

func buildKeyPairFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cvm.Filter {
	var filters []*cvm.Filter

	if equalsQuals["key_name"] != nil {
		if value := equalsQuals["key_name"].GetStringValue(); value != "" {
			filters = append(filters, &cvm.Filter{
				Name:   common.StringPtr("key-name"),
				Values: common.StringPtrs([]string{value}),
			})
		}
	}

	if equalsQuals["project_id"] != nil {
		value := equalsQuals["project_id"].GetInt64Value()
		filters = append(filters, &cvm.Filter{
			Name:   common.StringPtr("project-id"),
			Values: common.StringPtrs([]string{strconv.FormatInt(value, 10)}),
		})
	}

	return filters
}

// Get Function

func getCvmKeyPair(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	keyId := d.EqualsQuals["key_id"].GetStringValue()
	plugin.Logger(ctx).Info("getCvmKeyPair", "key_id", keyId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cvm.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCvmKeyPair", "client_init_error", err)
		return nil, err
	}

	if keyId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cvm.NewDescribeKeyPairsRequest()
	req.KeyIds = common.StringPtrs([]string{keyId})

	utils.LogRequest(ctx, "getCvmKeyPair", req)
	resp, err := client.DescribeKeyPairsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCvmKeyPair", "request_error", err)
		return nil, err
	}

	if resp == nil || resp.Response == nil || len(resp.Response.KeyPairSet) == 0 {
		return nil, nil
	}

	row := toCvmKeyPairRow(resp.Response.KeyPairSet[0])
	row.Region = region
	return row, nil
}

// Row Type

type cvmKeyPairRow struct {
	KeyId                 string
	KeyName               string
	ProjectId             int64
	Description           string
	PublicKey             string
	AssociatedInstanceIds []string
	CreatedTime           *time.Time
	Tags                  map[string]string
	Region                string
	Akas                  []string
}

func toCvmKeyPairRow(keyPair *cvm.KeyPair) cvmKeyPairRow {
	row := cvmKeyPairRow{
		KeyId:                 utils.PtrString(keyPair.KeyId),
		KeyName:               utils.PtrString(keyPair.KeyName),
		ProjectId:             utils.PtrInt64(keyPair.ProjectId),
		Description:           utils.PtrString(keyPair.Description),
		PublicKey:             utils.PtrString(keyPair.PublicKey),
		AssociatedInstanceIds: utils.PtrStringSlice(keyPair.AssociatedInstanceIds),
		CreatedTime:           utils.ParseTimestamp(keyPair.CreatedTime),
		Tags:                  TagsToMap(keyPair.Tags),
	}

	if row.KeyId != "" {
		row.Akas = []string{"tencentcloud:cvm:keypair:" + row.KeyId}
	}
	return row
}
