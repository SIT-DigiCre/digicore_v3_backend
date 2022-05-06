CREATE TABLE GroupUser
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id  BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    UNIQUE uk_user_id_group_id (user_id, group_id),
    CONSTRAINT fk_group_id_group_user FOREIGN KEY (group_id) REFERENCES `Group`(id),
    PRIMARY KEY (id)
);
