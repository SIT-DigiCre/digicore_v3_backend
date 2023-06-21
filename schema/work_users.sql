CREATE TABLE work_users
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    work_id    BINARY(16)   NOT NULL,
    user_id    BINARY(16)   NOT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_work_id_user_id (work_id, user_id),
    PRIMARY KEY (id)
);
