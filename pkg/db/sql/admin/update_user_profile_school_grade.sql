UPDATE user_profiles
INNER JOIN users ON users.id = user_profiles.user_id
LEFT JOIN (
    SELECT
        user_id,
        COALESCE(SUM(grade_diff), 0) AS approved_grade_diffs
    FROM grade_updates
    WHERE status = 'approved'
    GROUP BY user_id
) AS grade_update_summary ON grade_update_summary.user_id = users.id
SET user_profiles.school_grade =
    (
        YEAR(CURRENT_DATE)
        - IF(MONTH(CURRENT_DATE) BETWEEN 1 AND 3, 1, 0)
        - 2000
        - CAST(SUBSTRING(users.student_number, 3, 2) AS UNSIGNED)
        + 1
        + CASE LOWER(LEFT(users.student_number, 1))
            WHEN 'm' THEN 4
            WHEN 'n' THEN 6
            ELSE 0
          END
        + COALESCE(grade_update_summary.approved_grade_diffs, 0)
    )
WHERE user_profiles.is_graduated = false
  AND (
      user_profiles.is_member = true
      OR user_profiles.school_grade < 6
  )
  AND CHAR_LENGTH(users.student_number) >= 4
  AND SUBSTRING(users.student_number, 3, 2) REGEXP '^[0-9]{2}$';
