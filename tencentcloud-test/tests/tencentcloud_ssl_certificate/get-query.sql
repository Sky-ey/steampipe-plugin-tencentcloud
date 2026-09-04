select certificate_id, alias, domain
from tencentcloud_ssl_certificate
where certificate_id = 'cert-test'
order by certificate_id;
