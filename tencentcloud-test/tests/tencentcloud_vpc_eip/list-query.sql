select address_id, title, region
from tencentcloud_vpc_eip
where region = 'ap-guangzhou'
order by address_id;
