select rule_id, domain, action_type, ip, note
from tencentcloud_waf_ip_access_control
where domain = 'test.example.com'
  and action_type = 40
order by rule_id;
