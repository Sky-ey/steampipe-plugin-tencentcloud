select alarm_notice_id, name, type, region
from tencentcloud_cls_alarm_notice
where region = 'ap-guangzhou'
order by alarm_notice_id;
