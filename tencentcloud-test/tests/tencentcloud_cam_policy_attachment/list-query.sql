select policy_id, policy_name, entity_id, entity_name, related_type
from tencentcloud_cam_policy_attachment
order by policy_id, entity_id;
