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
    CASE WHEN gc.claim IS NOT NULL THEN true ELSE false END as is_admin,
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
    parent_address,
    COUNT(*) OVER() as total
FROM user_profiles
LEFT JOIN users ON users.id = user_profiles.user_id
LEFT JOIN user_private_profiles ON user_private_profiles.user_id = user_profiles.user_id
LEFT JOIN groups_users gu ON gu.user_id = user_profiles.user_id
LEFT JOIN group_claims gc ON gc.group_id = gu.group_id AND gc.claim = 'admin'
WHERE 1 = 1
/* IF query */
  AND (
    username LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
    OR users.student_number LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
  )
/* END */
/* IF schoolGrade */
  AND school_grade = /*schoolGrade*/0
/* END */
/* IF isAdmin */
  AND (CASE WHEN gc.claim IS NOT NULL THEN true ELSE false END) = /*isAdmin*/false
/* END */
GROUP BY user_profiles.user_id
ORDER BY user_profiles.created_at DESC
LIMIT /*limit*/100
OFFSET /*offset*/0;
