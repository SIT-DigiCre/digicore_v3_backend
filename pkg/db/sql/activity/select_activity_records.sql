SELECT
    BIN_TO_UUID(activities.id) AS record_id,
    BIN_TO_UUID(activities.user_id) AS user_id,
    user_profiles.username,
    activities.place,
    activities.checked_in_at,
    activities.checked_out_at,
    activities.initial_checked_in_at,
    activities.initial_checked_out_at
FROM activities
INNER JOIN user_profiles ON user_profiles.user_id = activities.user_id
ORDER BY activities.checked_in_at DESC
LIMIT /*limit*/50
/* IF offset */
OFFSET /*offset*/0
/* END */;
