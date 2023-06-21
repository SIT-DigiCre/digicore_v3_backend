SELECT BIN_TO_UUID(id) AS file_id, name FROM user_files where BIN_TO_UUID(id) in /*fileIds*/('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
