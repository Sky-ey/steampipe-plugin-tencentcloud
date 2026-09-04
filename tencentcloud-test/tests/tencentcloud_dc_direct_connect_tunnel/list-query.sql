select direct_connect_tunnel_id, direct_connect_tunnel_name, region
from tencentcloud_dc_direct_connect_tunnel
where region = 'ap-guangzhou'
order by direct_connect_tunnel_id;
