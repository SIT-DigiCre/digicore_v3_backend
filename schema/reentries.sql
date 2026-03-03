CREATE TABLE reentries
(
    id            BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id       BINARY(16)   NOT NULL,
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending',
    note          VARCHAR(255) NOT NULL DEFAULT '',
    checked_by    BINARY(16)            DEFAULT NULL,
    checked_at    DATETIME               DEFAULT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_reentries_user_id_created_at (user_id, created_at),
    INDEX idx_reentries_status_created_at (status, created_at)
);
