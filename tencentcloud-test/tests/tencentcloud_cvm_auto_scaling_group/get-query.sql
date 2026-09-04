select auto_scaling_group_id, title, region
from tencentcloud_cvm_auto_scaling_group
where auto_scaling_group_id = 'asg-test' and region = 'ap-guangzhou';
