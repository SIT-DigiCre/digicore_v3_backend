SELECT BIN_TO_UUID(users.id) AS user_id, users.student_number FROM users WHERE BIN_TO_UUID(users.id) IN /*userIds*/('00000000-0000-0000-0000-000000000000');
