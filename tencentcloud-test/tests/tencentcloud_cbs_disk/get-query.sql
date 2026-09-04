select disk_id, title, region
from tencentcloud_cbs_disk
where disk_id = 'disk-test'
  and region = 'ap-guangzhou';
