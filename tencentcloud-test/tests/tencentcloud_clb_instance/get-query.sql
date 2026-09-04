select load_balancer_id, title, region
from tencentcloud_clb_instance
where load_balancer_id = 'lb-test' and region = 'ap-guangzhou'
order by load_balancer_id;
