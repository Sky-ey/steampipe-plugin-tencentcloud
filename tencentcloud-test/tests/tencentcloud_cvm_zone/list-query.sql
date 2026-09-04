select zone, title, region
from tencentcloud_cvm_zone
where region = 'ap-guangzhou'
order by zone;
