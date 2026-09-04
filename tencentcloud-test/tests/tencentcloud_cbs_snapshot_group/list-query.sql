select snapshot_group_id, title, region
from tencentcloud_cbs_snapshot_group
where region = 'ap-guangzhou'
order by snapshot_group_id;
