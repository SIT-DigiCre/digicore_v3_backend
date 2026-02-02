SELECT
    COUNT(*) as total
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
/* END */;
