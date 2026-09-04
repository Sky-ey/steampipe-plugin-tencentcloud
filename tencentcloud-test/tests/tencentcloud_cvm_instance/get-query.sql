select instance_id, title, region
from tencentcloud_cvm_instance
where instance_id = 'ins-test'
  and region = 'ap-guangzhou';
