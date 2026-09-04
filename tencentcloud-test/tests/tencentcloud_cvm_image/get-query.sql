select image_id, title, region
from tencentcloud_cvm_image
where image_id = 'img-test'
  and region = 'ap-guangzhou';
