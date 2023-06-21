SELECT BIN_TO_UUID(user_id) AS user_id, username, icon_url FROM user_profiles where BIN_TO_UUID(id) in /*userIds*/('aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
