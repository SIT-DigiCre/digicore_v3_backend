SELECT BIN_TO_UUID(id) AS work_id, name FROM works /* IF authorId */ WHERE UUID_TO_BIN(/*authorId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') IN ( SELECT user_id FROM work_users WHERE work_id = works.id )  /* END */ ORDER BY updated_at DESC LIMIT 100 /* IF offset */ OFFSET /*offset*/0 /* END */;
