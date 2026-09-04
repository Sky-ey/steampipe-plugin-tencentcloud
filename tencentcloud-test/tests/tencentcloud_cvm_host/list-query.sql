select host_id, title, region
from tencentcloud_cvm_host
where region = 'ap-guangzhou'
order by host_id;
