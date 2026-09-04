select snapshot_group_id, title, region
from tencentcloud_cbs_snapshot_group
where snapshot_group_id = 'csnap-test'
  and region = 'ap-guangzhou';
