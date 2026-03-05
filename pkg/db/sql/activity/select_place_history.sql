SELECT
    BIN_TO_UUID(user_profiles.user_id) AS user_id,
    user_profiles.username,
    user_profiles.short_introduction,
    user_profiles.icon_url,
    /* IF periodIsDay */
    COUNT(activities.id) AS check_in_count
    /* ELSE */
    COUNT(DISTINCT DATE(activities.checked_in_at)) AS check_in_count
    /* END */
FROM user_profiles
INNER JOIN activities ON user_profiles.user_id = activities.user_id
WHERE activities.place = /*place*/'place'
    AND activities.checked_in_at >= /*startDate*/'2024-01-01 00:00:00'
    AND activities.checked_in_at <= /*endDate*/'2024-12-31 23:59:59'
GROUP BY user_profiles.user_id, user_profiles.username, user_profiles.short_introduction, user_profiles.icon_url
ORDER BY check_in_count DESC;

