SELECT BIN_TO_UUID(groups.id) as group_id, name, description, user_count, joinable, IF(1 <= count(user_id), true, false) as joined FROM `groups` LEFT JOIN groups_users ON groups_users.group_id = groups.id AND user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') GROUP BY groups.id LIMIT 50 /* IF offset*/ OFFSET /*offset*/0 /* END */;
