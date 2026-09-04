select target_group_id, title, region
from tencentcloud_clb_target_group
where target_group_id = 'tg-test' and region = 'ap-guangzhou'
order by target_group_id;
