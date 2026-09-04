select route_table_id, route_item_id, region
from tencentcloud_vpc_route
where region = 'ap-guangzhou'
order by route_item_id;
