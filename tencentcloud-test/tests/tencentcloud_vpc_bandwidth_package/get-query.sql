select bandwidth_package_id, title, region
from tencentcloud_vpc_bandwidth_package
where bandwidth_package_id = 'bp-test' and region = 'ap-guangzhou'
order by bandwidth_package_id;
