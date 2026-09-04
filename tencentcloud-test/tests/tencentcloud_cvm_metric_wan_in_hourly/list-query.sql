select instance_id, namespace, value, region
from tencentcloud_cvm_metric_wan_in_hourly
where region = 'ap-guangzhou'
order by instance_id;
