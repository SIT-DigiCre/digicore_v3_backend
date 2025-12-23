UPDATE activities
SET
    checked_out_at = /*checkedOutAt*/NOW()
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
