select ip, title, region
from tencentcloud_clb_cross_target
where region = 'ap-guangzhou'
order by ip;
