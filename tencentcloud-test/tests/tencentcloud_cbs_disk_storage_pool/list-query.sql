select cdc_id, title, region
from tencentcloud_cbs_disk_storage_pool
where region = 'ap-guangzhou'
order by cdc_id;
