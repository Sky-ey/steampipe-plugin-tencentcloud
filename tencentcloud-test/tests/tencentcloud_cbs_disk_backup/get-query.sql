select disk_backup_id, title, region
from tencentcloud_cbs_disk_backup
where disk_backup_id = 'dbp-test'
  and region = 'ap-guangzhou';
