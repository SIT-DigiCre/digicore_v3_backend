SELECT 
  BIN_TO_UUID(groups.id) as group_id,
  name,
  description,
  user_count,
  joinable,
  IF(1 <= count(user_id), true, false) as joined,
  IF(COUNT(group_claims.group_id) > 0, true, false) as is_admin_group
FROM `groups`
LEFT JOIN groups_users ON groups_users.group_id = groups.id AND user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
LEFT JOIN group_claims ON group_claims.group_id = groups.id AND group_claims.claim = 'admin'
WHERE groups.id = UUID_TO_BIN(/*groupId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
GROUP BY groups.id;
