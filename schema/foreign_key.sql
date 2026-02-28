ALTER TABLE event_reservation_users ADD CONSTRAINT fk_event_reservation_users_reservation_id_event_reservations_id FOREIGN KEY (reservation_id) REFERENCES event_reservations(id);
ALTER TABLE event_reservation_users ADD CONSTRAINT fk_event_reservation_users_reservation_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE event_reservations ADD CONSTRAINT fk_event_reservations_event_id_events_id FOREIGN KEY (event_id) REFERENCES events(id);

ALTER TABLE groups_users ADD CONSTRAINT fk_groups_users_group_id_groups_id FOREIGN KEY (group_id) REFERENCES `groups`(id);
ALTER TABLE groups_users ADD CONSTRAINT fk_groups_users_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_files ADD CONSTRAINT fk_user_files_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_payments ADD CONSTRAINT fk_user_payments_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_private_profiles ADD CONSTRAINT fk_user_private_profiles_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE activities ADD CONSTRAINT fk_activities_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE work_work_tags ADD CONSTRAINT fk_work_work_tags_tag_id_work_tags_id FOREIGN KEY (tag_id) REFERENCES work_tags(id);
ALTER TABLE work_work_tags ADD CONSTRAINT fk_work_work_tags_work_id_works_id FOREIGN KEY (work_id) REFERENCES works(id);

ALTER TABLE work_users ADD CONSTRAINT fk_work_users_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE work_users ADD CONSTRAINT fk_work_users_work_id_works_id FOREIGN KEY (work_id) REFERENCES works(id);

ALTER TABLE work_files ADD CONSTRAINT fk_work_files_file_id_user_files_id FOREIGN KEY (file_id) REFERENCES user_files(id);
ALTER TABLE work_files ADD CONSTRAINT fk_work_files_work_id_works_id FOREIGN KEY (work_id) REFERENCES works(id);

ALTER TABLE budgets ADD CONSTRAINT fk_budgets_proposer_user_id_users_id FOREIGN KEY (proposer_user_id) REFERENCES users(id);
ALTER TABLE budgets ADD CONSTRAINT fk_budgets_approver_user_id_users_id FOREIGN KEY (approver_user_id) REFERENCES users(id);

ALTER TABLE group_claims ADD CONSTRAINT fk_group_claims_group_id_groups_id FOREIGN KEY (group_id) REFERENCES `groups`(id);

ALTER TABLE user_profile_links ADD CONSTRAINT fk_user_profile_links_user_id_user_id  FOREIGN KEY(user_id)  REFERENCES users(id) ON DELETE CASCADE;