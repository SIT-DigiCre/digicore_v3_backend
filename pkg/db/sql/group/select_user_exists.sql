SELECT COUNT(*) > 0 AS user_exists
FROM users
WHERE id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
