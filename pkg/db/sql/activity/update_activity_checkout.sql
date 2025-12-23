UPDATE activities
SET
    initial_checked_out_at = /*checkedOutAt*/NOW(),
    checked_out_at = /*checkedOutAt*/NOW()
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
