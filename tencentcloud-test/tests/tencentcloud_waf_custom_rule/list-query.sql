select rule_id, rule_name, domain, status
from tencentcloud_waf_custom_rule
where domain = 'test.example.com'
order by rule_id;
