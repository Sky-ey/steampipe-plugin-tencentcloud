select topic_id, topic_name, logset_id, region
from tencentcloud_cls_topic
where region = 'ap-guangzhou'
order by topic_id;
