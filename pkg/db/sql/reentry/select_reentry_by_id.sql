SELECT BIN_TO_UUID(id) AS reentry_id, BIN_TO_UUID(user_id) AS user_id, status, note, created_at, updated_at FROM reentries WHERE id = UUID_TO_BIN(/*reentryId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
