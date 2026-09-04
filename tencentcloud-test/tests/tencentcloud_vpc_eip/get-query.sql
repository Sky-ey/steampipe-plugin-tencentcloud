select address_id, title, region
from tencentcloud_vpc_eip
where address_id = 'eip-test' and region = 'ap-guangzhou'
order by address_id;
