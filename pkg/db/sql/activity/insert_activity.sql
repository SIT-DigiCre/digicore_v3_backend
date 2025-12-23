INSERT INTO activities (
    id,
    user_id,
    place,
    note,
    initial_checked_in_at,
    initial_checked_out_at,
    checked_in_at,
    checked_out_at
) VALUES (
    UUID_TO_BIN(@id),
    UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee'),
    /*place*/'place',
    NULL,
    /*initialCheckedInAt*/NOW(),
    NULL,
    /*checkedInAt*/NOW(),
    NULL
);
