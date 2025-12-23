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
    /*initialCheckedInAt*/'1970-01-01 00:00:00',
    NULL,
    /*checkedInAt*/'1970-01-01 00:00:00',
    NULL
);
