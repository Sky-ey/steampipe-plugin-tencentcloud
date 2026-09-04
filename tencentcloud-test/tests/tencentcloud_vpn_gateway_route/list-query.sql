select vpn_gateway_id, route_id, destination_cidr_block, region
from tencentcloud_vpn_gateway_route
where region = 'ap-guangzhou'
order by route_id;
