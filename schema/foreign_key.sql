ALTER TABLE groups_users ADD CONSTRAINT fk_groups_users_group_id_groups_id FOREIGN KEY (group_id) REFERENCES `groups`(id);
ALTER TABLE groups_users ADD CONSTRAINT fk_groups_users_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_files ADD CONSTRAINT fk_user_files_user_id_users_id  FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_payments ADD CONSTRAINT fk_user_payments_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_private_profiles ADD CONSTRAINT fk_user_private_profiles_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE user_private_profiles ADD CONSTRAINT fk_user_profiles_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
