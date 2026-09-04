select snapshot_id, title, region
from tencentcloud_cbs_snapshot
where region = 'ap-guangzhou'
order by snapshot_id;
