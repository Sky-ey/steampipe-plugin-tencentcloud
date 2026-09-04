select direct_connect_tunnel_id, direct_connect_tunnel_name, region
from tencentcloud_dc_direct_connect_tunnel
where direct_connect_tunnel_id = 'dcx-test' and region = 'ap-guangzhou'
order by direct_connect_tunnel_id;
