select traffic_package_id, title, region
from tencentcloud_vpc_traffic_package
where traffic_package_id = 'tp-test' and region = 'ap-guangzhou'
order by traffic_package_id;
