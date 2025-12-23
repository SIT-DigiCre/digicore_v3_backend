UPDATE activities
SET
    checked_in_at = /*checkedInAt*/'1970-01-01 00:00:00',
    checked_out_at = /*checkedOutAt*/'1970-01-01 00:00:00'
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
