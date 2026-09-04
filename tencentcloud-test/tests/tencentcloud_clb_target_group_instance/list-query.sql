select target_group_id, instance_id, region
from tencentcloud_clb_target_group_instance
where target_group_id = 'tg-test' and region = 'ap-guangzhou'
order by instance_id;
