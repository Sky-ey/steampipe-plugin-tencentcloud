select launch_template_id, title, region
from tencentcloud_cvm_launch_template
where region = 'ap-guangzhou'
order by launch_template_id;
