select network_interface_id, title, region
from tencentcloud_vpc_eni
where network_interface_id = 'eni-test' and region = 'ap-guangzhou'
order by network_interface_id;
