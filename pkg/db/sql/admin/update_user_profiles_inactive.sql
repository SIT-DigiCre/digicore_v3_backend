UPDATE user_profiles
LEFT JOIN user_private_profiles
  ON user_private_profiles.user_id = user_profiles.user_id
SET user_profiles.is_member = false,
    user_private_profiles.first_name = '',
    user_private_profiles.last_name = '',
    user_private_profiles.first_name_kana = '',
    user_private_profiles.last_name_kana = '',
    user_private_profiles.is_male = true,
    user_private_profiles.phone_number = '',
    user_private_profiles.address = '',
    user_private_profiles.parent_name = '',
    user_private_profiles.parent_last_name = '',
    user_private_profiles.parent_first_name = '',
    user_private_profiles.parent_cellphone_number = '',
    user_private_profiles.parent_homephone_number = '',
    user_private_profiles.parent_address = ''
WHERE (user_profiles.active_limit < CURRENT_DATE
       AND user_profiles.is_member = true
       AND user_profiles.is_graduated = false)
   OR user_profiles.is_member = false;
-- is_member = falseも条件に入れることでこれまでの先輩方の個人情報も消えるよう調整
