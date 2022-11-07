CREATE TABLE groups_users
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id  BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    UNIQUE uk_user_id_group_id (user_id, group_id),
    PRIMARY KEY (id)
);
