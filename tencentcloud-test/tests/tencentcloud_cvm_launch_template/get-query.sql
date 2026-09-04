select launch_template_id, title, region
from tencentcloud_cvm_launch_template
where launch_template_id = 'lt-test'
  and region = 'ap-guangzhou';
