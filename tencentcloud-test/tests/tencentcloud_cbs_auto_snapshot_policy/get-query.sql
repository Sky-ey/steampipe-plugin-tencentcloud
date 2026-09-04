select auto_snapshot_policy_id, title, region
from tencentcloud_cbs_auto_snapshot_policy
where auto_snapshot_policy_id = 'asp-test'
  and region = 'ap-guangzhou';
