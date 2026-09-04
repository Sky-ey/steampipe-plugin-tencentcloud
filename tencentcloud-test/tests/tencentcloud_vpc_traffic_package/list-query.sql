select traffic_package_id, title, region
from tencentcloud_vpc_traffic_package
where region = 'ap-guangzhou'
order by traffic_package_id;
