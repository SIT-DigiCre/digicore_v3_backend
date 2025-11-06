UPDATE `groups`
SET user_count = user_count + 1, updated_at = NOW()
WHERE id = UUID_TO_BIN(/*groupId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
