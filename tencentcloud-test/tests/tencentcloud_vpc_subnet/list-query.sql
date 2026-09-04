select subnet_id, title, region
from tencentcloud_vpc_subnet
where region = 'ap-guangzhou'
order by subnet_id;
