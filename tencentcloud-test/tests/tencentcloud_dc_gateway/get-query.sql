select direct_connect_gateway_id, direct_connect_gateway_name, region
from tencentcloud_dc_gateway
where direct_connect_gateway_id = 'dcg-test' and region = 'ap-guangzhou'
order by direct_connect_gateway_id;
