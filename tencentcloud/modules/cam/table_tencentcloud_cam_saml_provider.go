package cam

import (
	"context"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
)

// Table Definition

func TableTencentcloudCamSamlProvider() *plugin.Table {
	return &plugin.Table{
		Name:        "tencentcloud_cam_saml_provider",
		Description: "Tencent Cloud CAM SAML identity providers used for enterprise SSO.",
		List: &plugin.ListConfig{
			Hydrate: listCamSamlProviders,
			Tags:    map[string]string{"service": "cam", "action": "ListSAMLProviders"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getCamSamlProvider,
			Tags:       map[string]string{"service": "cam", "action": "GetSAMLProvider"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getCamSamlProvider,
				Tags: map[string]string{"service": "cam", "action": "GetSAMLProvider"},
			},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the SAML identity provider.", Transform: transform.FromField("Name")},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the SAML identity provider.", Transform: transform.FromField("Description")},
			{Name: "create_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the SAML identity provider was created.", Transform: transform.FromField("CreateTime").NullIfZero()},
			{Name: "modify_time", Type: proto.ColumnType_TIMESTAMP, Description: "The time the SAML identity provider was last modified.", Transform: transform.FromField("ModifyTime").NullIfZero()},
			{Name: "saml_metadata", Type: proto.ColumnType_STRING, Description: "The SAML metadata document of the provider, only available from the GetSAMLProvider detail call.", Hydrate: getCamSamlProvider, Transform: transform.FromField("SAMLMetadata")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listCamSamlProviders(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listCamSamlProviders", "quals", d.EqualsQuals)

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listCamSamlProviders", "client_init_error", err)
		return nil, err
	}

	d.WaitForListRateLimit(ctx)
	req := cam.NewListSAMLProvidersRequest()
	utils.LogRequest(ctx, "listCamSamlProviders", req)
	resp, err := client.ListSAMLProvidersWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("listCamSamlProviders", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	for _, p := range resp.Response.SAMLProviderSet {
		if p == nil {
			continue
		}
		d.StreamListItem(ctx, toCamSamlProviderRow(p))
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

// Get Function

func getCamSamlProvider(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (any, error) {
	isSupplement := h.Item != nil
	var name string
	if isSupplement {
		if row, ok := h.Item.(camSamlProviderRow); ok {
			name = row.Name
		}
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}
	plugin.Logger(ctx).Debug("getCamSamlProvider", "name", name, "supplement", isSupplement)

	if name == "" {
		return nil, nil
	}

	client := &cam.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("getCamSamlProvider", "client_init_error", err)
		return nil, err
	}

	req := cam.NewGetSAMLProviderRequest()
	req.Name = common.StringPtr(name)

	utils.LogRequest(ctx, "getCamSamlProvider", req)
	resp, err := client.GetSAMLProviderWithContext(ctx, req)
	if err != nil {
		plugin.Logger(ctx).Error("getCamSamlProvider", "request_error", err)
		return nil, err
	}
	if resp == nil || resp.Response == nil {
		return nil, nil
	}

	// Supplemental hydrate: only the SAMLMetadata column consumes this value.
	if isSupplement {
		return resp.Response, nil
	}

	// Get hydrate: build a row from the response.
	row := camSamlProviderRow{
		Name:        utils.PtrString(resp.Response.Name),
		Description: utils.PtrString(resp.Response.Description),
		CreateTime:  utils.ParseTimestamp(resp.Response.CreateTime),
		ModifyTime:  utils.ParseTimestamp(resp.Response.ModifyTime),
	}
	if row.Name != "" {
		row.Title = row.Name
		row.Akas = []string{"tencentcloud:cam:saml-provider:" + row.Name}
	}
	return row, nil
}

// Row Type

type camSamlProviderRow struct {
	Name        string
	Description string
	CreateTime  *time.Time
	ModifyTime  *time.Time
	Title       string
	Akas        []string
}

func toCamSamlProviderRow(p *cam.SAMLProviderInfo) camSamlProviderRow {
	row := camSamlProviderRow{
		Name:        utils.PtrString(p.Name),
		Description: utils.PtrString(p.Description),
		CreateTime:  utils.ParseTimestamp(p.CreateTime),
		ModifyTime:  utils.ParseTimestamp(p.ModifyTime),
	}

	if row.Name != "" {
		row.Title = row.Name
		row.Akas = []string{"tencentcloud:cam:saml-provider:" + row.Name}
	}
	return row
}
