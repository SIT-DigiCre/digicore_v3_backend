UPDATE activities
SET
    checked_in_at = /*checkedInAt*/NOW(),
    checked_out_at = /*checkedOutAt*/NOW()
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
