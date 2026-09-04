select config_id, config_name, log_type, region
from tencentcloud_cls_config
where config_id = 'config-test'
  and region = 'ap-guangzhou';
