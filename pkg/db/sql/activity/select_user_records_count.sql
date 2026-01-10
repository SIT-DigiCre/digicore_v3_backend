SELECT
    COUNT(*) AS total
FROM activities
WHERE activities.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
/* IF place */
    AND activities.place = /*place*/'place'
/* END */;

