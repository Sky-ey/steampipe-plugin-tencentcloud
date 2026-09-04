select auto_snapshot_policy_id, title, region
from tencentcloud_cbs_auto_snapshot_policy
where region = 'ap-guangzhou'
order by auto_snapshot_policy_id;
