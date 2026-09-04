select uconfig_id, title, region
from tencentcloud_clb_customized_config
where config_type = 'CLB' and region = 'ap-guangzhou'
order by uconfig_id;
