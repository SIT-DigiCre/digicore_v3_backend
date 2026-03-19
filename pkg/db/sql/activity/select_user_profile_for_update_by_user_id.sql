SELECT
    BIN_TO_UUID(user_id) AS user_id
FROM
    user_profiles
WHERE
    user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
FOR UPDATE;
