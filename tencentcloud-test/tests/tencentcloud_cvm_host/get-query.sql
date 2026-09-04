select host_id, title, region
from tencentcloud_cvm_host
where host_id = 'host-test' and region = 'ap-guangzhou'
order by host_id;
