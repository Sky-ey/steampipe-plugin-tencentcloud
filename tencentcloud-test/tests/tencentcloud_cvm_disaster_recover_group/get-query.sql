select disaster_recover_group_id, name, title, region
from tencentcloud_cvm_disaster_recover_group
where disaster_recover_group_id = 'drg-test' and region = 'ap-guangzhou';
