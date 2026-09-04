select vpc_id, title, region
from tencentcloud_vpc
where region = 'ap-guangzhou'
order by vpc_id;
