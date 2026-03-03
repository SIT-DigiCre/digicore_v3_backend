UPDATE user_profiles
SET is_member = false
WHERE active_limit < CURRENT_DATE
  AND is_member = true;
