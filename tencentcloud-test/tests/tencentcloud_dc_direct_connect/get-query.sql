select direct_connect_id, direct_connect_name, region
from tencentcloud_dc_direct_connect
where direct_connect_id = 'dc-test' and region = 'ap-guangzhou'
order by direct_connect_id;
