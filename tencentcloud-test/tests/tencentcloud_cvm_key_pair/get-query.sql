select key_id, title, region
from tencentcloud_cvm_key_pair
where key_id = 'skey-test'
  and region = 'ap-guangzhou';
