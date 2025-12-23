UPDATE activities
SET
    checked_out_at = /*checkedOutAt*/'1970-01-01 00:00:00',
    note = COALESCE(/*note*/NULL, note)
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
