SELECT
    BIN_TO_UUID(user_profiles.user_id) AS user_id,
    user_profiles.username,
    user_profiles.short_introduction,
    user_profiles.icon_url,
    latest_activities.checked_in_at
FROM user_profiles
INNER JOIN (
    SELECT
        user_id,
        checked_in_at
    FROM (
        SELECT
            user_id,
            checked_in_at,
            ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY checked_in_at DESC) AS rn
        FROM activities
        WHERE place = /*place*/'place'
    ) AS ranked
    WHERE rn = 1
) AS latest_per_user ON user_profiles.user_id = latest_per_user.user_id
INNER JOIN activities AS latest_activities ON latest_per_user.user_id = latest_activities.user_id
    AND latest_per_user.checked_in_at = latest_activities.checked_in_at
    AND latest_activities.place = /*place*/'place'
    AND latest_activities.checked_out_at IS NULL
ORDER BY latest_activities.checked_in_at ASC;

