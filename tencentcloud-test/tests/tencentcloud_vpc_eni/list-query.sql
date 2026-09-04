select network_interface_id, title, region
from tencentcloud_vpc_eni
where region = 'ap-guangzhou'
order by network_interface_id;
