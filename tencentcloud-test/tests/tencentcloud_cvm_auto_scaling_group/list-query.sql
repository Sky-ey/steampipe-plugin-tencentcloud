select auto_scaling_group_id, title, region
from tencentcloud_cvm_auto_scaling_group
where region = 'ap-guangzhou'
order by auto_scaling_group_id;
