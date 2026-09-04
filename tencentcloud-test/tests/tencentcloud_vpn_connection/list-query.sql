select vpn_connection_id, vpn_connection_name, region
from tencentcloud_vpn_connection
where region = 'ap-guangzhou'
order by vpn_connection_id;
