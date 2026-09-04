select load_balancer_id, ip, region
from tencentcloud_clb_block_ip
where load_balancer_id = 'lb-test' and region = 'ap-guangzhou'
order by ip;
