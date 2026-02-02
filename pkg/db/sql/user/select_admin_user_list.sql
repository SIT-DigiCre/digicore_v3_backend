SELECT
    BIN_TO_UUID(user_profiles.user_id) as user_id,
    users.student_number,
    username,
    school_grade,
    icon_url,
    discord_userid,
    active_limit,
    short_introduction,
    introduction,
    IF(
        EXISTS(
            SELECT 1
            FROM groups_users
            INNER JOIN group_claims ON groups_users.group_id = group_claims.group_id
            WHERE groups_users.user_id = user_profiles.user_id
            AND group_claims.claim = 'admin'
        ),
        true,
        false
    ) as is_admin,
    first_name,
    last_name,
    first_name_kana,
    last_name_kana,
    is_male,
    phone_number,
    address,
    parent_name,
    parent_last_name,
    parent_first_name,
    parent_cellphone_number,
    parent_homephone_number,
    parent_address
FROM user_profiles
LEFT JOIN users ON users.id = user_profiles.user_id
LEFT JOIN user_private_profiles ON user_private_profiles.user_id = user_profiles.user_id
WHERE 1 = 1
/* IF query */
  AND (
    username LIKE CONCAT('%', /*query*/'', '%')
    OR users.student_number LIKE CONCAT('%', /*query*/'', '%')
  )
/* END */
/* IF schoolGrade */
  AND school_grade = /*schoolGrade*/0
/* END */
/* IF isAdmin */
  AND IF(
        EXISTS(
            SELECT 1
            FROM groups_users
            INNER JOIN group_claims ON groups_users.group_id = group_claims.group_id
            WHERE groups_users.user_id = user_profiles.user_id
            AND group_claims.claim = 'admin'
        ),
        true,
        false
    ) = /*isAdmin*/false
/* END */
ORDER BY user_profiles.created_at DESC
LIMIT /*limit*/100
OFFSET /*offset*/0;
