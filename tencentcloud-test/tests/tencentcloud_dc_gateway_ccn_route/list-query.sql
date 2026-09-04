select direct_connect_gateway_id, route_id, destination_cidr_block, region
from tencentcloud_dc_gateway_ccn_route
where region = 'ap-guangzhou'
order by route_id;
