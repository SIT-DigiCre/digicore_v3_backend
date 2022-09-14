-- +migrate Up
ALTER TABLE `user_profiles` ADD CONSTRAINT `fk_user_profiles_user_id_users_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`);
ALTER TABLE `user_private_profiles` ADD CONSTRAINT `fk_user_private_profiles_user_id_users_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`);
ALTER TABLE `user_payments` ADD CONSTRAINT `fk_user_payments_user_id_users_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`);
ALTER TABLE `event_reservations` ADD CONSTRAINT `fk_event_reservations_event_id_events_id` FOREIGN KEY (`event_id`) REFERENCES `events`(`id`);
ALTER TABLE `event_reservation_users` ADD CONSTRAINT `fk_event_reservation_users_reservation_id_event_reservations_id` FOREIGN KEY (`reservation_id`) REFERENCES `event_reservations`(`id`);
ALTER TABLE `event_reservation_users` ADD CONSTRAINT `fk_event_reservation_users_reservation_id_users_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`);

-- +migrate Down
