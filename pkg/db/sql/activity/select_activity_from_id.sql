SELECT
    BIN_TO_UUID(id) AS id,
    BIN_TO_UUID(user_id) AS user_id,
    place,
    note,
    initial_checked_in_at,
    initial_checked_out_at,
    checked_in_at,
    checked_out_at,
    created_at,
    updated_at
FROM
    activities
WHERE
    id = UUID_TO_BIN(/*id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
