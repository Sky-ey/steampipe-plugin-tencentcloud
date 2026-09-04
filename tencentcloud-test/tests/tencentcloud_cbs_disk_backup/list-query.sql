select disk_backup_id, title, region
from tencentcloud_cbs_disk_backup
where region = 'ap-guangzhou'
order by disk_backup_id;
