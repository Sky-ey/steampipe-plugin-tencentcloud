package cls

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cls "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cls/v20201016"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudClsTopic() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_cls_topic",
		Description:       "Tencent Cloud CLS topics are the ingestion and storage units of Cloud Log Service; each topic belongs to a logset and holds the log data, indexing, lifecycle and storage configuration for a given log stream.",
		GetMatrixItemFunc: buildClsRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "topic_id", Require: plugin.Optional},
				{Name: "topic_name", Require: plugin.Optional},
				{Name: "logset_id", Require: plugin.Optional},
			},
			Hydrate: listClsTopics,
			Tags:    map[string]string{"service": "cls", "action": "DescribeTopics"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("topic_id"),
			Hydrate:    getClsTopic,
			Tags:       map[string]string{"service": "cls", "action": "DescribeTopics"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "topic_id", Type: proto.ColumnType_STRING, Description: "The unique ID of the log topic.", Transform: transform.FromField("TopicId")},
			{Name: "topic_name", Type: proto.ColumnType_STRING, Description: "The name of the log topic.", Transform: transform.FromField("TopicName")},
			{Name: "logset_id", Type: proto.ColumnType_STRING, Description: "The ID of the logset the topic belongs to.", Transform: transform.FromField("LogsetId")},
			{Name: "biz_type", Type: proto.ColumnType_INT, Description: "The topic type (0: log topic, 1: metric topic).", Transform: transform.FromField("BizType")},
			{Name: "partition_count", Type: proto.ColumnType_INT, Description: "The number of partitions in the topic.", Transform: transform.FromField("PartitionCount")},
			{Name: "storage_type", Type: proto.ColumnType_STRING, Description: "The storage type of the topic (hot: standard storage, cold: infrequent storage).", Transform: transform.FromField("StorageType")},
			{Name: "period", Type: proto.ColumnType_INT, Description: "The lifecycle in days (1-3600); a value of 3640 indicates permanent retention.", Transform: transform.FromField("Period")},
			{Name: "hot_period", Type: proto.ColumnType_INT, Description: "The standard storage lifecycle in days when log settlement is enabled (hotPeriod < Period); 0 means log settlement is not enabled.", Transform: transform.FromField("HotPeriod")},
			{Name: "index_enabled", Type: proto.ColumnType_BOOL, Description: "Whether the topic has indexing enabled (log topics only).", Transform: transform.FromField("Index")},
			{Name: "status", Type: proto.ColumnType_BOOL, Description: "Whether log collection is enabled (true: enabled, false: disabled).", Transform: transform.FromField("Status")},
			{Name: "auto_split", Type: proto.ColumnType_BOOL, Description: "Whether auto-split is enabled for the topic.", Transform: transform.FromField("AutoSplit")},
			{Name: "max_split_partitions", Type: proto.ColumnType_INT, Description: "The maximum number of partitions allowed when auto-split is enabled.", Transform: transform.FromField("MaxSplitPartitions")},
			{Name: "is_web_tracking", Type: proto.ColumnType_BOOL, Description: "Whether free (anonymous) access is enabled for the log topic.", Transform: transform.FromField("IsWebTracking")},
			{Name: "is_source_from", Type: proto.ColumnType_BOOL, Description: "Whether to record public network source IP and server receipt time.", Transform: transform.FromField("IsSourceFrom")},
			{Name: "billing_mode", Type: proto.ColumnType_INT, Description: "The current billing mode (0: function billing by usage, 1: billing by raw log size).", Transform: transform.FromField("BillingMode")},
			{Name: "new_billing_mode", Type: proto.ColumnType_INT, Description: "The billing mode that takes effect after an async task succeeds (0: function billing by usage, 1: billing by raw log size).", Transform: transform.FromField("NewBillingMode")},
			{Name: "migration_status", Type: proto.ColumnType_INT, Description: "The async migration status (1: in progress, 2: completed, 3: failure, 4: canceled).", Transform: transform.FromField("MigrationStatus")},
			{Name: "topic_async_task_id", Type: proto.ColumnType_STRING, Description: "The async migration task ID, if any.", Transform: transform.FromField("TopicAsyncTaskID")},
			{Name: "effective_date", Type: proto.ColumnType_TIMESTAMP, Description: "The expected effective date after async migration.", Transform: transform.FromField("EffectiveDate").NullIfZero()},
			{Name: "key_id", Type: proto.ColumnType_STRING, Description: "The kms-cls service key ID used by the topic.", Transform: transform.FromField("KeyId")},
			{Name: "assumer_uin", Type: proto.ColumnType_INT, Description: "The Uin of the service provider that created the topic, if any.", Transform: transform.FromField("AssumerUin")},
			{Name: "assumer_name", Type: proto.ColumnType_STRING, Description: "The cloud product identifier that created the topic (e.g., CDN, TKE), if any.", Transform: transform.FromField("AssumerName")},
			{Name: "sub_assumer_name", Type: proto.ColumnType_STRING, Description: "The cloud product sub-identifier (e.g., TKE-Audit, TKE-Event), if any.", Transform: transform.FromField("SubAssumerName")},
			{Name: "role_name", Type: proto.ColumnType_STRING, Description: "The role used by the service provider that created the topic, if any.", Transform: transform.FromField("RoleName")},
			{Name: "describes", Type: proto.ColumnType_STRING, Description: "The description of the topic.", Transform: transform.FromField("Describes")},
			{Name: "extends", Type: proto.ColumnType_JSON, Description: "The extended information of the log topic.", Transform: transform.FromField("Extends")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the topic.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the topic as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the topic resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("TopicName")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listClsTopics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listClsTopics", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listClsTopics", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	filters := buildClsTopicFilters(d.EqualsQuals)

	pageSize := utils.PageSizeFromLimit(d, 100, 1)
	var offset int64

	for {
		d.WaitForListRateLimit(ctx)

		req := cls.NewDescribeTopicsRequest()
		req.Offset = &offset
		req.Limit = &pageSize
		if len(filters) > 0 {
			req.Filters = filters
		}

		utils.LogRequest(ctx, "listClsTopics", req)
		resp, err := client.DescribeTopicsWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listClsTopics", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		items := resp.Response.Topics
		for _, item := range items {
			if item == nil {
				continue
			}
			row := toClsTopicRow(item)
			row.Region = region
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.Response == nil {
			break
		}
		total := utils.PtrInt64(resp.Response.TotalCount)
		if utils.PageDone(offset, int64(len(items)), total) {
			break
		}
		offset += int64(len(items))
	}

	return nil, nil
}

func buildClsTopicFilters(equalsQuals plugin.KeyColumnEqualsQualMap) []*cls.Filter {
	return utils.BuildFilters(equalsQuals, newClsFilter,
		utils.FilterMapping{QualName: "topic_id", FilterName: "topicId"},
		utils.FilterMapping{QualName: "topic_name", FilterName: "topicName"},
		utils.FilterMapping{QualName: "logset_id", FilterName: "logsetId"},
	)
}

// Get Function

func getClsTopic(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	topicId := d.EqualsQuals["topic_id"].GetStringValue()
	plugin.Logger(ctx).Debug("getClsTopic", "topic_id", topicId, "region", d.EqualsQualString(utils.MatrixKeyRegion))

	client := &cls.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getClsTopic", "client_init_error", err)
		return nil, err
	}

	if topicId == "" {
		return nil, nil
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	req := cls.NewDescribeTopicsRequest()
	req.Filters = []*cls.Filter{newClsFilter("topicId", topicId)}

	utils.LogRequest(ctx, "getClsTopic", req)
	resp, err := client.DescribeTopicsWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getClsTopic", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Topics) == 0 {
		return nil, nil
	}

	row := toClsTopicRow(resp.Response.Topics[0])
	row.Region = region
	return row, nil
}

// Row Type

type clsTopicRow struct {
	TopicId            string
	TopicName          string
	LogsetId           string
	BizType            uint64
	PartitionCount     int64
	StorageType        string
	Period             int64
	HotPeriod          uint64
	Index              bool
	Status             bool
	AutoSplit          bool
	MaxSplitPartitions int64
	IsWebTracking      bool
	IsSourceFrom       bool
	BillingMode        uint64
	NewBillingMode     uint64
	MigrationStatus    uint64
	TopicAsyncTaskID   string
	EffectiveDate      *time.Time
	KeyId              string
	AssumerUin         uint64
	AssumerName        string
	SubAssumerName     string
	RoleName           string
	Describes          string
	Extends            *cls.TopicExtendInfo
	CreateTime         *time.Time
	Tags               map[string]string
	Region             string
	Akas               []string
}

func toClsTopicRow(t *cls.TopicInfo) clsTopicRow {
	row := clsTopicRow{
		TopicId:            utils.PtrString(t.TopicId),
		TopicName:          utils.PtrString(t.TopicName),
		LogsetId:           utils.PtrString(t.LogsetId),
		BizType:            utils.PtrUint64(t.BizType),
		PartitionCount:     utils.PtrInt64(t.PartitionCount),
		StorageType:        utils.PtrString(t.StorageType),
		Period:             utils.PtrInt64(t.Period),
		HotPeriod:          utils.PtrUint64(t.HotPeriod),
		Index:              utils.PtrBool(t.Index),
		Status:             utils.PtrBool(t.Status),
		AutoSplit:          utils.PtrBool(t.AutoSplit),
		MaxSplitPartitions: utils.PtrInt64(t.MaxSplitPartitions),
		IsWebTracking:      utils.PtrBool(t.IsWebTracking),
		IsSourceFrom:       utils.PtrBool(t.IsSourceFrom),
		BillingMode:        utils.PtrUint64(t.BillingMode),
		NewBillingMode:     utils.PtrUint64(t.NewBillingMode),
		MigrationStatus:    utils.PtrUint64(t.MigrationStatus),
		TopicAsyncTaskID:   utils.PtrString(t.TopicAsyncTaskID),
		EffectiveDate:      utils.ParseTimestamp(t.EffectiveDate),
		KeyId:              utils.PtrString(t.KeyId),
		AssumerUin:         utils.PtrUint64(t.AssumerUin),
		AssumerName:        utils.PtrString(t.AssumerName),
		SubAssumerName:     utils.PtrString(t.SubAssumerName),
		RoleName:           utils.PtrString(t.RoleName),
		Describes:          utils.PtrString(t.Describes),
		Extends:            t.Extends,
		CreateTime:         utils.ParseTimestamp(t.CreateTime),
		Tags:               TagsToMap(t.Tags),
	}

	if row.TopicId != "" {
		row.Akas = []string{"tencentcloud:cls:topic:" + row.TopicId}
	}
	return row
}
