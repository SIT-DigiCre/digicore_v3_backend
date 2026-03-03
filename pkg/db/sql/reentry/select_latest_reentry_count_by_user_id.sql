SELECT COALESCE(MAX(reentry_count), 0) AS reentry_count FROM reentries WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
