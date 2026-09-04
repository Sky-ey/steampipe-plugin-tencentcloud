select cluster_id, cluster_name, cluster_type, cluster_status, region
from tencentcloud_tke_cluster
where region = 'ap-guangzhou'
order by cluster_id;
