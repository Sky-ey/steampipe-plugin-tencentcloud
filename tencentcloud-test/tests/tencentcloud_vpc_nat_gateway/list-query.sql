select nat_gateway_id, title, region
from tencentcloud_vpc_nat_gateway
where region = 'ap-guangzhou'
order by nat_gateway_id;
