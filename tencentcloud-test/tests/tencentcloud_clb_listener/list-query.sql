select listener_id, title, region
from tencentcloud_clb_listener
where load_balancer_id = 'lb-test' and region = 'ap-guangzhou'
order by listener_id;
