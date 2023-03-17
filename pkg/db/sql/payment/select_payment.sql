SELECT BIN_TO_UUID(user_payments.id) AS id, BIN_TO_UUID(user_payments.user_id) AS user_id, transfer_name, checked, student_number FROM user_payments LEFT JOIN users ON users.id = user_payments.user_id WHERE `year` = /*year*/2022 LIMIT 50 /* IF offset*/ OFFSET /*offset*/0 /* END */;
