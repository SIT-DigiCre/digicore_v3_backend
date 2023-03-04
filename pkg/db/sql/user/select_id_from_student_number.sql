SELECT BIN_TO_UUID(id) AS id, IF(CURRENT_TIMESTAMP <= active_limit,true,false)  AS active FROM users WHERE student_number = /*studentNumber*/'aa21000';
