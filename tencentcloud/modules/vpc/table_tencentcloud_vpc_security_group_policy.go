package vpc

import (
	"context"
	"github.com/tencentcloud/steampipe-plugin-tencentcloud/tencentcloud/utils"
	"strconv"
	"time"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/vpc/v20170312"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Table Definition

func TableTencentcloudVpcSecurityGroupPolicy() *plugin.Table {
	return &plugin.Table{
		Name:              "tencentcloud_vpc_security_group_policy",
		Description:       "Tencent Cloud VPC security group rules (ingress and egress), expanded one row per policy.",
		GetMatrixItemFunc: buildVpcRegionList,
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "security_group_id", Require: plugin.Optional},
				{Name: "direction", Require: plugin.Optional},
			},
			Hydrate: listVpcSecurityGroupPolicies,
			Tags:    map[string]string{"service": "vpc", "action": "DescribeSecurityGroupPolicies"},
		},
		Columns: utils.WithCommonColumns([]*plugin.Column{
			{Name: "security_group_id", Type: proto.ColumnType_STRING, Description: "The security group ID that owns the rule.", Transform: transform.FromField("SecurityGroupId")},
			{Name: "direction", Type: proto.ColumnType_STRING, Description: "The direction of the rule (INBOUND, OUTBOUND).", Transform: transform.FromField("Direction")},
			{Name: "policy_index", Type: proto.ColumnType_INT, Description: "The index number of the rule within the security group.", Transform: transform.FromField("PolicyIndex")},
			{Name: "protocol", Type: proto.ColumnType_STRING, Description: "The protocol (TCP, UDP, ICMP, ICMPv6, ALL).", Transform: transform.FromField("Protocol")},
			{Name: "port", Type: proto.ColumnType_STRING, Description: "The port range (`all`, a single port, or a port range).", Transform: transform.FromField("Port")},
			{Name: "cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv4 CIDR block the rule applies to.", Transform: transform.FromField("CidrBlock")},
			{Name: "ipv6_cidr_block", Type: proto.ColumnType_STRING, Description: "The IPv6 CIDR block the rule applies to.", Transform: transform.FromField("Ipv6CidrBlock")},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "The action to take (ACCEPT, DROP).", Transform: transform.FromField("Action")},
			{Name: "policy_description", Type: proto.ColumnType_STRING, Description: "The description of the rule.", Transform: transform.FromField("PolicyDescription")},
			{Name: "referenced_security_group_id", Type: proto.ColumnType_STRING, Description: "The peer security group ID referenced by the rule, if any.", Transform: transform.FromField("ReferencedSecurityGroupId")},
			{Name: "service_template", Type: proto.ColumnType_JSON, Description: "The protocol port (template) specification referenced by the rule.", Transform: transform.FromField("ServiceTemplate")},
			{Name: "address_template", Type: proto.ColumnType_JSON, Description: "The IP address (template) specification referenced by the rule.", Transform: transform.FromField("AddressTemplate")},
			{Name: "modify_time", Type: proto.ColumnType_TIMESTAMP, Description: "The last modification time of the rule.", Transform: transform.FromField("ModifyTime").NullIfZero()},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The region where the rule resides.", Transform: transform.FromField("Region")},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource.", Transform: transform.FromField("Title")},
			{Name: "akas", Type: proto.ColumnType_JSON, Description: "Array of globally unique identifier strings (also known as) for the resource.", Transform: transform.FromField("Akas")},
		}),
	}
}

// List Function

func listVpcSecurityGroupPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("listVpcSecurityGroupPolicies", "region", d.EqualsQualString(utils.MatrixKeyRegion), "quals", d.EqualsQuals)

	client := &vpc.Client{}
	if err := utils.InitClient(ctx, d, client); err != nil {
		plugin.Logger(ctx).Error("listVpcSecurityGroupPolicies", "client_init_error", err)
		return nil, err
	}

	region := d.EqualsQualString(utils.MatrixKeyRegion)

	directionFilter := d.EqualsQualString("direction")

	sgIds, err := resolveVpcSecurityGroupIDs(ctx, d, client)
	if err != nil {
		plugin.Logger(ctx).Error("listVpcSecurityGroupPolicies", "resolve_security_group_ids_error", err)
		return nil, err
	}

	for _, sgId := range sgIds {
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		req := vpc.NewDescribeSecurityGroupPoliciesRequest()
		req.SecurityGroupId = &sgId

		utils.LogRequest(ctx, "listVpcSecurityGroupPolicies", req)
		resp, err := client.DescribeSecurityGroupPoliciesWithContext(ctx, req)
		if err != nil {
			plugin.Logger(ctx).Error("listVpcSecurityGroupPolicies", "request_error", err)
			return nil, err
		}
		if resp == nil || resp.Response == nil || resp.Response.SecurityGroupPolicySet == nil {
			continue
		}

		policySet := resp.Response.SecurityGroupPolicySet

		if directionFilter == "" || directionFilter == "INBOUND" {
			for _, p := range policySet.Ingress {
				if p == nil {
					continue
				}
				row := toVpcSecurityGroupPolicyRow(p, sgId, "INBOUND", region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
		if directionFilter == "" || directionFilter == "OUTBOUND" {
			for _, p := range policySet.Egress {
				if p == nil {
					continue
				}
				row := toVpcSecurityGroupPolicyRow(p, sgId, "OUTBOUND", region)
				d.StreamListItem(ctx, row)
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

// Row Type

type vpcSecurityGroupPolicyRow struct {
	SecurityGroupId           string
	Direction                 string
	PolicyIndex               int64
	Protocol                  string
	Port                      string
	CidrBlock                 string
	Ipv6CidrBlock             string
	Action                    string
	PolicyDescription         string
	ReferencedSecurityGroupId string
	ServiceTemplate           *vpc.ServiceTemplateSpecification
	AddressTemplate           *vpc.AddressTemplateSpecification
	ModifyTime                *time.Time
	Region                    string
	Title                     string
	Akas                      []string
}

func toVpcSecurityGroupPolicyRow(p *vpc.SecurityGroupPolicy, sgId, direction, region string) vpcSecurityGroupPolicyRow {
	row := vpcSecurityGroupPolicyRow{
		SecurityGroupId:           sgId,
		Direction:                 direction,
		PolicyIndex:               utils.PtrInt64(p.PolicyIndex),
		Protocol:                  utils.PtrString(p.Protocol),
		Port:                      utils.PtrString(p.Port),
		CidrBlock:                 utils.PtrString(p.CidrBlock),
		Ipv6CidrBlock:             utils.PtrString(p.Ipv6CidrBlock),
		Action:                    utils.PtrString(p.Action),
		PolicyDescription:         utils.PtrString(p.PolicyDescription),
		ReferencedSecurityGroupId: utils.PtrString(p.SecurityGroupId),
		ServiceTemplate:           p.ServiceTemplate,
		AddressTemplate:           p.AddressTemplate,
		ModifyTime:                utils.ParseTimestamp(p.ModifyTime),
		Region:                    region,
	}

	row.Title = row.Protocol + "/" + row.Port
	if row.CidrBlock != "" {
		row.Title += " -> " + row.CidrBlock
	} else if row.Ipv6CidrBlock != "" {
		row.Title += " -> " + row.Ipv6CidrBlock
	} else if row.ReferencedSecurityGroupId != "" {
		row.Title += " -> " + row.ReferencedSecurityGroupId
	}

	row.Akas = []string{"tencentcloud:vpc:security-group-policy:" + sgId + "/" + direction + "/" + strconv.FormatInt(row.PolicyIndex, 10)}
	return row
}
