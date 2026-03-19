UPDATE user_profiles
SET is_graduated = true
WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
