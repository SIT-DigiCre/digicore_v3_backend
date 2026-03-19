UPDATE user_profiles
INNER JOIN (
    SELECT
        UUID_TO_BIN(update_rows.user_id) AS user_id,
        update_rows.school_grade AS school_grade
    FROM JSON_TABLE(
        /*updatesJson*/'[]',
        '$[*]' COLUMNS (
            user_id CHAR(36) PATH '$.userId',
            school_grade INT PATH '$.schoolGrade'
        )
    ) AS update_rows
) AS updates ON updates.user_id = user_profiles.user_id
SET user_profiles.school_grade = updates.school_grade
WHERE user_profiles.is_graduated = false
  AND (
      user_profiles.is_member = true
      OR user_profiles.school_grade < 6
  );
