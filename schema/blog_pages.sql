CREATE TABLE blog_pages
(
    id          BINARY(16)      NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id     BINARY(16)      NOT NULL,
    title       VARCHAR(255)    NOT NULL,
    body        TEXT            NOT NULL,
    created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
)