SELECT claim FROM groups_users RIGHT JOIN group_claims ON groups_users.group_id = group_claims.group_id WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
