select vpn_connection_id, vpn_connection_name, region
from tencentcloud_vpn_connection
where vpn_connection_id = 'vpnx-test' and region = 'ap-guangzhou'
order by vpn_connection_id;
