select alarm_notice_id, name, type, region
from tencentcloud_cls_alarm_notice
where alarm_notice_id = 'notice-test'
  and region = 'ap-guangzhou';
