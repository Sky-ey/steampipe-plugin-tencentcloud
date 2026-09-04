select config_id, config_name, topic_id, log_type, region
from tencentcloud_cls_config
where region = 'ap-guangzhou'
order by config_id;
