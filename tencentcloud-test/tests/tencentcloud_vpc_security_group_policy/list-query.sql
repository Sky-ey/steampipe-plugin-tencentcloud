select security_group_id, direction, protocol, region
from tencentcloud_vpc_security_group_policy
where security_group_id = 'sg-test' and region = 'ap-guangzhou'
order by direction, protocol;
