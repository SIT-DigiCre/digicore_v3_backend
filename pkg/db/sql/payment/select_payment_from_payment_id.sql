SELECT BIN_TO_UUID(user_payments.id) AS id, BIN_TO_UUID(user_payments.user_id) AS user_id, transfer_name, checked, student_number FROM user_payments LEFT JOIN users ON users.id = user_payments.user_id WHERE user_payments.id = UUID_TO_BIN(/*paymentId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
