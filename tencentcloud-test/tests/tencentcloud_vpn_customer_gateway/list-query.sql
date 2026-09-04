select customer_gateway_id, customer_gateway_name, region
from tencentcloud_vpn_customer_gateway
where region = 'ap-guangzhou'
order by customer_gateway_id;
