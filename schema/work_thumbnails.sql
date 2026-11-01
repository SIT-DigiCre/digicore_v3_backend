CREATE TABLE work_thumbnails
(
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id BINARY(16) NOT NULL,
    storage_key VARCHAR(255) NOT NULL,
    content_type VARCHAR(32) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
