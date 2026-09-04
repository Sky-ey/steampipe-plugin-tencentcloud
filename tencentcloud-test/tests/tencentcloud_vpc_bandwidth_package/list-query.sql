select bandwidth_package_id, title, region
from tencentcloud_vpc_bandwidth_package
where region = 'ap-guangzhou'
order by bandwidth_package_id;
