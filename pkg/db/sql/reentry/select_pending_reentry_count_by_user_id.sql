SELECT COUNT(*) AS pending_count FROM reentries WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') AND status = 'pending';
