package cos

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudCosObject() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cos_object",
		Description: "Objects stored in Tencent Cloud Object Storage (COS) buckets.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Optional},
				{Name: "region", Require: plugin.Optional},
				{Name: "prefix", Require: plugin.Optional},
				{Name: "delimiter", Require: plugin.Optional},
			},
			Hydrate: listCosObjects,
			Tags:    map[string]string{"service": "cos", "action": "GetBucket"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "key", Type: proto.ColumnType_STRING, Description: "The object key (name).", Transform: transform.FromField("Key")},
			{Name: "bucket_name", Type: proto.ColumnType_STRING, Description: "The bucket that contains the object.", Transform: transform.FromField("BucketName")},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the bucket resides.", Transform: transform.FromField("Region")},
			{Name: "size", Type: proto.ColumnType_INT, Description: "The object size in bytes.", Transform: transform.FromField("Size")},
			{Name: "etag", Type: proto.ColumnType_STRING, Description: "The ETag (entity tag) of the object.", Transform: transform.FromField("ETag")},
			{Name: "storage_class", Type: proto.ColumnType_STRING, Description: "The storage class of the object (STANDARD, STANDARD_IA, ARCHIVE, DEEP_ARCHIVE, INTELLIGENT_TIERING).", Transform: transform.FromField("StorageClass")},
			{Name: "storage_tier", Type: proto.ColumnType_STRING, Description: "The storage tier for INTELLIGENT_TIERING objects (FREQUENT, INFREQUENT).", Transform: transform.FromField("StorageTier")},
			{Name: "last_modified", Type: proto.ColumnType_TIMESTAMP, Description: "The last modification time of the object.", Transform: transform.FromField("LastModified").NullIfZero()},
			{Name: "restore_status", Type: proto.ColumnType_STRING, Description: "The restore status for archived objects.", Transform: transform.FromField("RestoreStatus")},
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "The owner information of the object.", Transform: transform.FromField("Owner")},
			{Name: "prefix", Type: proto.ColumnType_STRING, Description: "The prefix used to filter objects; populated when specified in the query.", Transform: transform.FromQual("prefix")},
			{Name: "delimiter", Type: proto.ColumnType_STRING, Description: "The delimiter used to group keys; populated when specified in the query.", Transform: transform.FromQual("delimiter")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Key")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCosObjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCosObjects", "quals", d.EqualsQuals)

	bucketNames, regions, err := resolveCosBucketNames(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCosObjects", "resolve_bucket_names_error", err)
		return nil, err
	}

	prefix := d.EqualsQualString("prefix")
	delimiter := d.EqualsQualString("delimiter")

	for i, bucketName := range bucketNames {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
		region := regions[i]

		client, err := buildCosClient(ctx, d, bucketName, region)
		if err != nil {
			plugin.Logger(ctx).Error("listCosObjects", "client_init_error", err, "bucket", bucketName)
			return nil, err
		}

		maxKeys := int(maxKeysFromLimit(d))
		var marker string
		for {
			d.WaitForListRateLimit(ctx)

			opt := &cos.BucketGetOptions{
				Prefix:    prefix,
				Delimiter: delimiter,
				Marker:    marker,
				MaxKeys:   maxKeys,
			}

			utils.LogRequest(ctx, "listCosObjects", opt)
			res, _, err := client.Bucket.Get(ctx, opt)
			if err != nil {
				plugin.Logger(ctx).Error("listCosObjects", "request_error", err, "bucket", bucketName)
				return nil, err
			}
			if res == nil {
				break
			}

			for _, obj := range res.Contents {
				row := toCosObjectRow(obj, bucketName, region)
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
	}

	return nil, nil
}

// Row Type

type cosObjectRow struct {
	Key           string
	BucketName    string
	Region        string
	Size          int64
	ETag          string
	StorageClass  string
	StorageTier   string
	LastModified  *time.Time
	RestoreStatus string
	Owner         *cos.Owner
	Akas          []string
}

func toCosObjectRow(o cos.Object, bucketName, region string) cosObjectRow {
	row := cosObjectRow{
		Key:           o.Key,
		BucketName:    bucketName,
		Region:        region,
		Size:          o.Size,
		ETag:          o.ETag,
		StorageClass:  o.StorageClass,
		StorageTier:   o.StorageTier,
		RestoreStatus: o.RestoreStatus,
		Owner:         o.Owner,
	}

	row.LastModified = utils.ParseTimestamp(&o.LastModified)

	if row.Key != "" && row.BucketName != "" {
		row.Akas = []string{"tencentcloud:cos:object:" + row.Region + "/" + row.BucketName + "/" + row.Key}
	}
	return row
}
