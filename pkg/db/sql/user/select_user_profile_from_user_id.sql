SELECT
    BIN_TO_UUID(user_profiles.user_id) as user_id,
    student_number,
    username,
    school_grade,
    icon_url,
    discord_userid,
    active_limit,
    short_introduction,
    IF(EXISTS(
        SELECT 1
        FROM groups_users
        INNER JOIN group_claims ON groups_users.group_id = group_claims.group_id
        WHERE groups_users.user_id = user_profiles.user_id
        AND group_claims.claim IN /*adminClaims*/('account', 'infra')
    ), true, false) as is_admin
FROM user_profiles
LEFT JOIN users ON users.id = user_profiles.user_id
WHERE user_profiles.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
