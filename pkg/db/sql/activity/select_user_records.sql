SELECT
    BIN_TO_UUID(activities.id) AS record_id,
    activities.place,
    activities.checked_in_at,
    activities.checked_out_at,
    activities.initial_checked_in_at,
    activities.initial_checked_out_at
FROM activities
WHERE activities.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
/* IF place */
    AND activities.place = /*place*/'place'
/* END */
ORDER BY activities.checked_in_at DESC
LIMIT /*limit*/50
/* IF offset */
OFFSET /*offset*/0
/* END */;

