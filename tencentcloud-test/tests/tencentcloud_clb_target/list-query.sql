select load_balancer_id, listener_id, instance_id, region
from tencentcloud_clb_target
where load_balancer_id = 'lb-test' and region = 'ap-guangzhou'
order by instance_id;
