UPDATE activities
SET
    initial_checked_out_at = COALESCE(initial_checked_out_at, /*checkedOutAt*/'1970-01-01 00:00:00'),
    checked_out_at = /*checkedOutAt*/'1970-01-01 00:00:00',
    note = COALESCE(/*note*/NULL, note)
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
    AND checked_out_at IS NULL;
