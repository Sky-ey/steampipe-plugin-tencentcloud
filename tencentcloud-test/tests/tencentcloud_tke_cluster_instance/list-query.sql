select instance_id, cluster_id, instance_role, instance_state, region
from tencentcloud_tke_cluster_instance
where region = 'ap-guangzhou'
  and cluster_id = 'cls-test'
order by instance_id;
