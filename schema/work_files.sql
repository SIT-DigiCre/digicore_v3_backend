CREATE TABLE work_files
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    work_id     BINARY(16)   NOT NULL,
    file_id     BINARY(16)   NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (work_id, file_id),
    PRIMARY KEY (id)
);
