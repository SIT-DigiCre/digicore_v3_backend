DELETE user_private_profiles
FROM user_private_profiles
INNER JOIN user_profiles ON user_private_profiles.user_id = user_profiles.user_id
WHERE user_profiles.active_limit <= CURRENT_DATE - INTERVAL 1 YEAR AND (user_profiles.is_member = false OR user_profiles.is_graduated=true);
