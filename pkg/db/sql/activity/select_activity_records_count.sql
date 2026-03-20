SELECT COUNT(*) AS total
FROM activities
INNER JOIN user_profiles ON user_profiles.user_id = activities.user_id;
