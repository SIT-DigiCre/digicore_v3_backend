CREATE TABLE groups_users
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id  BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    UNIQUE uk_user_id_group_id (user_id, group_id),
    CONSTRAINT fk_groups_users_group_id_groups_id FOREIGN KEY (group_id) REFERENCES `groups`(id),
    CONSTRAINT fk_groups_users_user_id_users_id FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (id)
);
