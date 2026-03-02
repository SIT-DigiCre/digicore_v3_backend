CREATE TABLE grade_updates (
    id           BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id      BINARY(16) NOT NULL,
    grade_diff   INT NOT NULL DEFAULT -1,
    reason       TEXT NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'pending',
    approved_by  BINARY(16) NULL,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
