select cluster_id, cluster_name, cluster_status, region
from tencentcloud_tke_cluster
where cluster_id = 'cls-test'
  and region = 'ap-guangzhou';
