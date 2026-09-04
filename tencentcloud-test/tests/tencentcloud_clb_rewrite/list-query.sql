select load_balancer_id, listener_id, region
from tencentcloud_clb_rewrite
where load_balancer_id = 'lb-test' and region = 'ap-guangzhou'
order by listener_id;
