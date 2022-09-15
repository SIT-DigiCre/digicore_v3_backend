CREATE TABLE work_users
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    work_id    BINARY(16)   NOT NULL,
    user_id    BINARY(16)   NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (work_id, user_id),
    PRIMARY KEY (id)
);
