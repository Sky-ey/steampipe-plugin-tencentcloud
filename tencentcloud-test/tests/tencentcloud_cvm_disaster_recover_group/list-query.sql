select disaster_recover_group_id, name, title, region
from tencentcloud_cvm_disaster_recover_group
where region = 'ap-guangzhou'
order by disaster_recover_group_id;
