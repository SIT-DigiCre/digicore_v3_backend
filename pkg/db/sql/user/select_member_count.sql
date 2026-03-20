SELECT COUNT(*) as count 
FROM user_profiles 
WHERE active_limit >= CURDATE();
