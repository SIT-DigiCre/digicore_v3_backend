SELECT
  BIN_TO_UUID(users.id) AS user_id,
  users.student_number,
  COALESCE(SUM(CASE WHEN grade_updates.status = 'approved' THEN grade_updates.grade_diff ELSE 0 END), 0) AS approved_grade_diffs
FROM users
INNER JOIN user_profiles ON user_profiles.user_id = users.id
LEFT JOIN grade_updates ON grade_updates.user_id = users.id
WHERE user_profiles.is_graduated = false
  AND (
    user_profiles.is_member = true
    OR user_profiles.school_grade < 6
  )
GROUP BY users.id, users.student_number;
