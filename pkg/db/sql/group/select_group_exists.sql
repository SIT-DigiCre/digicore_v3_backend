SELECT COUNT(*) > 0 AS group_exists
FROM `groups`
WHERE id = UUID_TO_BIN(/*groupId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
