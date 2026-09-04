select alarm_id, alarm_name, status, region
from tencentcloud_cls_alarm
where region = 'ap-guangzhou'
order by alarm_id;
