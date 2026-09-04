select alarm_id, alarm_name, status, region
from tencentcloud_cls_alarm
where alarm_id = 'alarm-test'
  and region = 'ap-guangzhou';
