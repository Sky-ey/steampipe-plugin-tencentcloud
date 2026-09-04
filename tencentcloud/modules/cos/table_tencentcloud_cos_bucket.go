package cos

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCosBucket() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cos_bucket",
		Description: "Tencent Cloud Object Storage (COS) buckets under the current account.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "region", Require: plugin.Optional},
				{Name: "tag_key", Require: plugin.Optional},
				{Name: "tag_value", Require: plugin.Optional},
			},
			Hydrate: listCosBuckets,
			Tags:    map[string]string{"service": "cos", "action": "GetService"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The bucket name in the format <name>-<appid>.", Transform: transform.FromField("Name")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the bucket resides (e.g., ap-guangzhou).", Transform: transform.FromField("Region")},
			{Name: "creation_date", Type: proto.ColumnType_TIMESTAMP, Description: "The creation time of the bucket.", Transform: transform.FromField("CreationDate").NullIfZero()},
			{Name: "bucket_type", Type: proto.ColumnType_STRING, Description: "The bucket type (e.g., cos, cos_intl).", Transform: transform.FromField("BucketType")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type attribute of the bucket.", Transform: transform.FromField("Type")},
			{Name: "tag_key", Type: proto.ColumnType_STRING, Description: "The tag key used to filter buckets; populated when specified in the query.", Transform: transform.FromQual("tag_key")},
			{Name: "tag_value", Type: proto.ColumnType_STRING, Description: "The tag value used to filter buckets; populated when specified in the query.", Transform: transform.FromQual("tag_value")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Name")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCosBuckets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCosBuckets", "quals", d.EqualsQuals)

	client, err := buildCosClient(ctx, d, "", "")
	if err != nil {
		plugin.Logger(ctx).Error("listCosBuckets", "client_init_error", err)
		return nil, err
	}

	nameQual := d.EqualsQualString("name")
	regionQual := d.EqualsQualString("region")
	tagKey := d.EqualsQualString("tag_key")
	tagValue := d.EqualsQualString("tag_value")

	maxKeys := maxKeysFromLimit(d)
	var marker string
	for {
		d.WaitForListRateLimit(ctx)

		opt := &cos.ServiceGetOptions{
			Marker:   marker,
			MaxKeys:  maxKeys,
			Region:   regionQual,
			TagKey:   tagKey,
			TagValue: tagValue,
		}

		utils.LogRequest(ctx, "listCosBuckets", opt)
		res, _, err := client.Service.Get(ctx, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listCosBuckets", "request_error", err)
			return nil, err
		}
		if res == nil {
			break
		}

		for _, b := range res.Buckets {
			row := toCosBucketRow(b)
			if nameQual != "" && row.Name != nameQual {
				continue
			}
			d.StreamListItem(ctx, row)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !res.IsTruncated || res.NextMarker == "" {
			break
		}
		marker = res.NextMarker
	}

	return nil, nil
}

// Row Type

type cosBucketRow struct {
	Name         string
	Region       string
	CreationDate *time.Time
	BucketType   string
	Type         string
	Akas         []string
}

func toCosBucketRow(b cos.Bucket) cosBucketRow {
	row := cosBucketRow{
		Name:       b.Name,
		Region:     b.Region,
		BucketType: b.BucketType,
		Type:       b.Type,
	}

	row.CreationDate = utils.ParseTimestamp(&b.CreationDate)

	if row.Name != "" && row.Region != "" {
		row.Akas = []string{"tencentcloud:cos:bucket:" + row.Region + "/" + row.Name}
	}
	return row
}
