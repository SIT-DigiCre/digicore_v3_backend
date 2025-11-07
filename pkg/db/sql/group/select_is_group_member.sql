SELECT COUNT(*) > 0 AS is_member
FROM groups_users
WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  AND group_id = UUID_TO_BIN(/*groupId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
