SELECT BIN_TO_UUID(`budgets`.id) as budget_id, BIN_TO_UUID(proposer_user_id) as user_id, username as user_name, icon_url, `status`, class, name, budget, settlement, `budgets`.updated_at as updated_at FROM `budgets` LEFT JOIN `user_profiles` ON `budgets`.proposer_user_id = `user_profiles`.user_id;