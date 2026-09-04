select topic_id, topic_name, logset_id, region
from tencentcloud_cls_topic
where topic_id = 'topic-test'
  and region = 'ap-guangzhou';
