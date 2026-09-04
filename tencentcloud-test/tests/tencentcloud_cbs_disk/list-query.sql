select disk_id, title, region
from tencentcloud_cbs_disk
where region = 'ap-guangzhou'
order by disk_id;
