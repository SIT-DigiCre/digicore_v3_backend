SELECT BIN_TO_UUID(id) AS tag_id, name, description FROM work_tags WHERE id = UUID_TO_BIN(/*tagId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
