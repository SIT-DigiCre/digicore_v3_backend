CREATE TABLE GroupUser
(
    id       VARCHAR(36) NOT NULL,
    user_id  BINARY(16)  NOT NULL,
    group_id BINARY(16)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (user_id, group_id)
);
