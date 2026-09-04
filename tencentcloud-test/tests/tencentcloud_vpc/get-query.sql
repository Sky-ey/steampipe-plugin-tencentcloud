select vpc_id, title, region
from tencentcloud_vpc
where vpc_id = 'vpc-test' and region = 'ap-guangzhou'
order by vpc_id;
