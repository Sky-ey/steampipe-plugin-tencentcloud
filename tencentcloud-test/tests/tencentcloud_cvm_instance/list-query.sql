select instance_id, title, region
from tencentcloud_cvm_instance
where region = 'ap-guangzhou'
order by instance_id;
