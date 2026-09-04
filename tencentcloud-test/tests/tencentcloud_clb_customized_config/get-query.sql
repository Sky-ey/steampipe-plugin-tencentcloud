select uconfig_id, title, region
from tencentcloud_clb_customized_config
where uconfig_id = 'uconfig-test' and region = 'ap-guangzhou'
order by uconfig_id;
