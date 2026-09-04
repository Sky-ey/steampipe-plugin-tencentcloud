select instance_id, metric_name, namespace, value, region
from tencentcloud_cvm_metric_cpu_hourly
where region = 'ap-guangzhou'
order by instance_id;
