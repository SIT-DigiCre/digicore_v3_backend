SELECT first_name, last_name, first_name_kana, last_name_kana, is_male, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address FROM user_private_profiles WHERE user_id = UUID_TO_BIN(/*userID*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
