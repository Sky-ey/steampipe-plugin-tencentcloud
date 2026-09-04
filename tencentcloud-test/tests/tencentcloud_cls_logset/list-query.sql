select logset_id, logset_name, topic_count, region
from tencentcloud_cls_logset
where region = 'ap-guangzhou'
order by logset_id;
