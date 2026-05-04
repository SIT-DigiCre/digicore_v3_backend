SELECT
    BIN_TO_UUID(users.id) AS user_id,
    users.student_number,
    COALESCE(grade_update_summary.approved_grade_diffs, 0) AS approved_grade_diffs
FROM user_profiles
INNER JOIN users ON users.id = user_profiles.user_id
LEFT JOIN (
    SELECT
        user_id,
        COALESCE(SUM(grade_diff), 0) AS approved_grade_diffs
    FROM grade_updates
    WHERE status = 'approved'
    GROUP BY user_id
) AS grade_update_summary ON grade_update_summary.user_id = users.id
WHERE user_profiles.is_graduated = false
  AND (
      user_profiles.is_member = true
      OR user_profiles.school_grade < 6
  )
  AND (
      (
          LOWER(LEFT(users.student_number, 1)) IN ('m', 'n')
          AND CHAR_LENGTH(users.student_number) >= 4
          AND SUBSTRING(users.student_number, 2, 2) REGEXP '^[0-9]{2}$'
      )
      OR (
          LOWER(LEFT(users.student_number, 1)) NOT IN ('m', 'n')
          AND CHAR_LENGTH(users.student_number) >= 4
          AND SUBSTRING(users.student_number, 3, 2) REGEXP '^[0-9]{2}$'
      )
  );
