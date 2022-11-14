SELECT BIN_TO_UUID(id) AS work_id, name, description FROM works WHERE id = UUID_TO_BIN(/*workId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
