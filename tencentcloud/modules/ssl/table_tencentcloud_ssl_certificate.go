package ssl

import (
	"context"
	"strconv"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/ssl/v20191205"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudSslCertificate() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_ssl_certificate",
		Description: "Tencent Cloud SSL certificates managed under the SSL Certificates Service (document/api/400).",
		// SSL is an account-scoped global service: no GetMatrixItemFunc, no region column.
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "certificate_id", Require: plugin.Optional},
				{Name: "certificate_type", Require: plugin.Optional},
				{Name: "search_key", Require: plugin.Optional},
				{Name: "project_id", Require: plugin.Optional},
				{Name: "deployable", Require: plugin.Optional},
				{Name: "is_expiring", Require: plugin.Optional},
			},
			Hydrate: listSslCertificates,
			Tags:    map[string]string{"service": "ssl", "action": "DescribeCertificates"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("certificate_id"),
			Hydrate:    getSslCertificate,
			Tags:       map[string]string{"service": "ssl", "action": "DescribeCertificates"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "certificate_id", Type: proto.ColumnType_STRING, Description: "The certificate instance ID.", Transform: transform.FromField("CertificateId")},
			{Name: "alias", Type: proto.ColumnType_STRING, Description: "The remark name (alias) of the certificate.", Transform: transform.FromField("Alias")},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The primary domain name bound to the certificate.", Transform: transform.FromField("Domain")},
			{Name: "status", Type: proto.ColumnType_INT, Description: "The numeric certificate status (0 = under review, 1 = approved, 3 = expired, 7 = canceled, 10 = revoked, etc.).", Transform: transform.FromField("Status")},
			{Name: "status_name", Type: proto.ColumnType_STRING, Description: "The human-readable certificate status name.", Transform: transform.FromField("StatusName")},
			{Name: "certificate_type", Type: proto.ColumnType_STRING, Description: "The certificate type: `CA` for client certificate, `SVR` for server certificate.", Transform: transform.FromField("CertificateType")},
			{Name: "search_key", Type: proto.ColumnType_STRING, Description: "The search keyword used to filter certificates; populated when specified in the query.", Transform: transform.FromQual("search_key")},
			{Name: "cert_begin_time", Type: proto.ColumnType_TIMESTAMP, Description: "The certificate validity start time.", Transform: transform.FromField("CertBeginTime")},
			{Name: "cert_end_time", Type: proto.ColumnType_TIMESTAMP, Description: "The certificate expiration time.", Transform: transform.FromField("CertEndTime")},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The project ID to which the certificate belongs.", Transform: transform.FromField("ProjectId")},
			{Name: "project_info", Type: proto.ColumnType_JSON, Description: "The project information (name, creator UIN, creation time, etc.) associated with the certificate.", Transform: transform.FromField("ProjectInfo")},
			{Name: "subject_alt_name", Type: proto.ColumnType_JSON, Description: "Multiple domain names contained in the certificate (including the primary domain name).", Transform: transform.FromField("SubjectAltName")},
			{Name: "cert_sans", Type: proto.ColumnType_JSON, Description: "Multiple domain names (SANs) associated with the certificate.", Transform: transform.FromField("CertSANs")},
			{Name: "bound_resource", Type: proto.ColumnType_JSON, Description: "The cloud resources bound to the certificate.", Transform: transform.FromField("BoundResource")},
			{Name: "hosting_resource_types", Type: proto.ColumnType_JSON, Description: "The hosted resource type list for the certificate.", Transform: transform.FromField("HostingResourceTypes")},
			{Name: "deployable", Type: proto.ColumnType_BOOL, Description: "Whether the certificate can be deployed.", Transform: transform.FromField("Deployable")},
			{Name: "is_expiring", Type: proto.ColumnType_BOOL, Description: "Whether the certificate is about to expire (within 30 days).", Transform: transform.FromField("IsExpiring")},
			{Name: "renew_able", Type: proto.ColumnType_BOOL, Description: "Whether the certificate is renewable.", Transform: transform.FromField("RenewAble")},
			{Name: "is_dv", Type: proto.ColumnType_BOOL, Description: "Whether it is a DV (domain validation) certificate.", Transform: transform.FromField("IsDv")},
			{Name: "is_wildcard", Type: proto.ColumnType_BOOL, Description: "Whether it is a wildcard domain name certificate.", Transform: transform.FromField("IsWildcard")},
			{Name: "is_sm", Type: proto.ColumnType_BOOL, Description: "Whether it is a China SM (ShangMi) certificate.", Transform: transform.FromField("IsSM")},
			{Name: "is_vip", Type: proto.ColumnType_BOOL, Description: "Whether the customer is a VIP customer.", Transform: transform.FromField("IsVip")},
			{Name: "allow_download", Type: proto.ColumnType_BOOL, Description: "Whether the certificate is allowed to be downloaded.", Transform: transform.FromField("AllowDownload")},
			{Name: "validity_period", Type: proto.ColumnType_STRING, Description: "The certificate validity period in months.", Transform: transform.FromField("ValidityPeriod")},
			{Name: "insert_time", Type: proto.ColumnType_TIMESTAMP, Description: "The certificate creation time.", Transform: transform.FromField("InsertTime")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "The tags attached to the certificate as a map of key -> value.", Transform: transform.FromField("Tags")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listSslCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listSslCertificates", "quals", d.EqualsQuals)

	client := &ssl.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listSslCertificates", "client_init_error", err)
		return nil, err
	}

	pageSize := uint64(utils.PageSizeFromLimit(d, 1000, 1))
	offset := uint64(0)

	certIds := utils.SingleIdFromQual(d.EqualsQuals, "certificate_id")

	for {
		d.WaitForListRateLimit(ctx)

		req := ssl.NewDescribeCertificatesRequest()
		req.Limit = &pageSize
		req.Offset = &offset

		// Mode B: certificate_id takes priority over other quals.
		if len(certIds) > 0 {
			req.CertIds = certIds
		} else {
			applySslCertificateQuals(req, d.EqualsQuals)
		}

		utils.LogRequest(ctx, "listSslCertificates", req)
		resp, err := client.DescribeCertificatesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listSslCertificates", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil {
			break
		}

		certs := resp.Response.Certificates
		total := utils.PtrUint64(resp.Response.TotalCount)

		for _, item := range certs {
			if item == nil {
				continue
			}
			d.StreamListItem(ctx, toSslCertificateRow(item))
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if utils.PageDone(int64(offset), int64(len(certs)), total) {
			break
		}
		offset += uint64(len(certs))
	}

	return nil, nil
}

func applySslCertificateQuals(req *ssl.DescribeCertificatesRequest, equalsQuals plugin.KeyColumnEqualsQualMap) {
	if v := equalsQuals["certificate_type"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			req.CertificateType = &s
		}
	}
	if v := equalsQuals["search_key"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			req.SearchKey = &s
		}
	}
	if v := equalsQuals["project_id"]; v != nil {
		if s := v.GetStringValue(); s != "" {
			if id, err := strconv.ParseUint(s, 10, 64); err == nil {
				req.ProjectId = &id
			}
		}
	}
	if v := equalsQuals["deployable"]; v != nil {
		if v.GetBoolValue() {
			req.Deployable = common.Uint64Ptr(1)
		} else {
			req.Deployable = common.Uint64Ptr(0)
		}
	}
	if v := equalsQuals["is_expiring"]; v != nil {
		if v.GetBoolValue() {
			req.FilterExpiring = common.Uint64Ptr(1)
		} else {
			req.FilterExpiring = common.Uint64Ptr(0)
		}
	}
}

// Get Function

func getSslCertificate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	certId := d.EqualsQuals["certificate_id"].GetStringValue()
	plugin.Logger(ctx).Info("getSslCertificate", "certificate_id", certId)

	client := &ssl.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getSslCertificate", "client_init_error", err)
		return nil, err
	}

	if certId == "" {
		return nil, nil
	}

	req := ssl.NewDescribeCertificatesRequest()
	req.CertIds = common.StringPtrs([]string{certId})

	utils.LogRequest(ctx, "getSslCertificate", req)
	resp, err := client.DescribeCertificatesWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getSslCertificate", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil || len(resp.Response.Certificates) == 0 {
		return nil, nil
	}

	return toSslCertificateRow(resp.Response.Certificates[0]), nil
}

// Row Type

type sslCertificateRow struct {
	CertificateId        string
	Alias                string
	Domain               string
	Status               uint64
	StatusName           string
	CertificateType      string
	CertBeginTime        *time.Time
	CertEndTime          *time.Time
	ProjectId            string
	ProjectInfo          *ssl.ProjectInfo
	SubjectAltName       []string
	CertSANs             []string
	BoundResource        []string
	HostingResourceTypes []string
	Deployable           bool
	IsExpiring           bool
	RenewAble            bool
	IsDv                 bool
	IsWildcard           bool
	IsSM                 bool
	IsVip                bool
	AllowDownload        bool
	ValidityPeriod       string
	InsertTime           *time.Time
	Tags                 map[string]string
	Title                string
	Akas                 []string
}

func toSslCertificateRow(c *ssl.Certificates) sslCertificateRow {
	row := sslCertificateRow{
		CertificateId:        utils.PtrString(c.CertificateId),
		Alias:                utils.PtrString(c.Alias),
		Domain:               utils.PtrString(c.Domain),
		Status:               utils.PtrUint64(c.Status),
		StatusName:           utils.PtrString(c.StatusName),
		CertificateType:      utils.PtrString(c.CertificateType),
		CertBeginTime:        utils.ParseTimestamp(c.CertBeginTime),
		CertEndTime:          utils.ParseTimestamp(c.CertEndTime),
		ProjectId:            utils.PtrString(c.ProjectId),
		ProjectInfo:          c.ProjectInfo,
		SubjectAltName:       utils.PtrStringSlice(c.SubjectAltName),
		CertSANs:             utils.PtrStringSlice(c.CertSANs),
		BoundResource:        utils.PtrStringSlice(c.BoundResource),
		HostingResourceTypes: utils.PtrStringSlice(c.HostingResourceTypes),
		Deployable:           utils.PtrBool(c.Deployable),
		IsExpiring:           utils.PtrBool(c.IsExpiring),
		RenewAble:            utils.PtrBool(c.RenewAble),
		IsDv:                 utils.PtrBool(c.IsDv),
		IsWildcard:           utils.PtrBool(c.IsWildcard),
		IsSM:                 utils.PtrBool(c.IsSM),
		IsVip:                utils.PtrBool(c.IsVip),
		AllowDownload:        utils.PtrBool(c.AllowDownload),
		ValidityPeriod:       utils.PtrString(c.ValidityPeriod),
		InsertTime:           utils.ParseTimestamp(c.InsertTime),
		Tags:                 TagsToMap(c.Tags),
	}

	if row.Alias != "" {
		row.Title = row.Alias
	} else {
		row.Title = row.CertificateId
	}

	if row.CertificateId != "" {
		row.Akas = []string{"tencentcloud:ssl:certificate:" + row.CertificateId}
	}
	return row
}
