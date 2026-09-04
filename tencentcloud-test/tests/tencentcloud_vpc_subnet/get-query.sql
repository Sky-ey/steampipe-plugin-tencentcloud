select subnet_id, title, region
from tencentcloud_vpc_subnet
where subnet_id = 'subnet-test' and region = 'ap-guangzhou'
order by subnet_id;
