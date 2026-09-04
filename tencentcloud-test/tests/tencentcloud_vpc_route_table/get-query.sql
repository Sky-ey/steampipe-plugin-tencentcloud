select route_table_id, title, region
from tencentcloud_vpc_route_table
where route_table_id = 'rtb-test' and region = 'ap-guangzhou'
order by route_table_id;
