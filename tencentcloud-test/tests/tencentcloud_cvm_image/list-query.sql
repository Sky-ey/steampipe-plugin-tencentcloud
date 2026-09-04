select image_id, title, region
from tencentcloud_cvm_image
where region = 'ap-guangzhou'
order by image_id;
