SELECT
    BIN_TO_UUID(users.id) AS id,
    COALESCE(user_profiles.is_member, false) AS is_member
FROM users
LEFT JOIN user_profiles ON users.id = user_profiles.user_id
WHERE student_number = /*studentNumber*/'aa21000';
