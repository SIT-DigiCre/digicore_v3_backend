SELECT BIN_TO_UUID(file_id) AS file_id, name FROM work_files LEFT JOIN user_files ON work_files.file_id = user_files.id WHERE work_id = UUID_TO_BIN(/*workId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
