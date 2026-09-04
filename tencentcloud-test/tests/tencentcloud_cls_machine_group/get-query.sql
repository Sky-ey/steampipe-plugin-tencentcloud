select machine_group_id, machine_group_name, os_type, region
from tencentcloud_cls_machine_group
where machine_group_id = 'mg-test'
  and region = 'ap-guangzhou';
