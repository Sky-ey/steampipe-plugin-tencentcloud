select key_id, title, region
from tencentcloud_cvm_key_pair
where region = 'ap-guangzhou'
order by key_id;
