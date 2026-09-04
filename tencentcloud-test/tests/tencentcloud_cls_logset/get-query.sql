select logset_id, logset_name, region
from tencentcloud_cls_logset
where logset_id = 'logset-test'
  and region = 'ap-guangzhou';
