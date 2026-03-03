SELECT
    COALESCE(is_member, false) AS is_member
FROM user_profiles
WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
