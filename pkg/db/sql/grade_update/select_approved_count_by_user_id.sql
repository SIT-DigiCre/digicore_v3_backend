SELECT COUNT(*) AS approved_count FROM grade_updates WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') AND status = 'approved';
