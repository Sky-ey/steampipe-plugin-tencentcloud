select load_balancer_id, title, region
from tencentcloud_clb_instance
where region = 'ap-guangzhou'
order by load_balancer_id;
