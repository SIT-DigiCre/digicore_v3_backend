CREATE TABLE work_work_tags
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    work_id    BINARY(16)   NOT NULL,
    tag_id     BINARY(16)   NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_work_id_tag_id (work_id, tag_id),
    PRIMARY KEY (id)
);
