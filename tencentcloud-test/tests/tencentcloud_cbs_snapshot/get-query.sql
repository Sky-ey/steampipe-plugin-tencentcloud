select snapshot_id, title, region
from tencentcloud_cbs_snapshot
where snapshot_id = 'snap-test'
  and region = 'ap-guangzhou';
