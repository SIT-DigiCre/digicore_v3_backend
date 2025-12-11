SELECT
    BIN_TO_UUID(user_id) as user_id,
    username,
    icon_url,
    short_introduction
FROM user_profiles
/* IF query */
WHERE username LIKE CONCAT('%', /*query*/'', '%') ESCAPE '\'
/* END */
LIMIT 100;
