select vpn_gateway_id, vpn_gateway_name, region
from tencentcloud_vpn_gateway
where region = 'ap-guangzhou'
order by vpn_gateway_id;
