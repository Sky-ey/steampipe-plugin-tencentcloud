select security_group_id, title, region
from tencentcloud_vpc_security_group
where security_group_id = 'sg-test' and region = 'ap-guangzhou'
order by security_group_id;
