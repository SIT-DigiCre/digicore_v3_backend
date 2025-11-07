SELECT COUNT(*) > 0 AS is_admin_group
FROM group_claims
WHERE group_id = UUID_TO_BIN(/*groupId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  AND claim = 'admin';
