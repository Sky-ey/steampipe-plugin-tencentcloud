select instance_id, namespace, value, region
from tencentcloud_cvm_metric_disk_hourly
where region = 'ap-guangzhou'
order by instance_id;
