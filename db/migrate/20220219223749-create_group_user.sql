-- +migrate Up
CREATE TABLE GroupUser
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id  BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE GroupUser;
