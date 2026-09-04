select
  i.instance_id,
  i.title as instance_title,
  d.disk_id,
  d.title as disk_title,
  i.region
from tencentcloud_cvm_instance as i
join tencentcloud_cbs_disk as d
  on d.instance_id = i.instance_id
 and d.region = i.region
where i.region = 'ap-guangzhou'
  and d.region = 'ap-guangzhou'
order by i.instance_id, d.disk_id;
