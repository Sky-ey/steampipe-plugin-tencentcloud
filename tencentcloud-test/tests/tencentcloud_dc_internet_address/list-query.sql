select instance_id, subnet, addr_type, region
from tencentcloud_dc_internet_address
where addr_type = 0 and region = 'ap-guangzhou'
order by instance_id;
