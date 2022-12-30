CREATE TABLE blogs
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id     BINARY(16)   NOT NULL,
    name        VARCHAR(255) NOT NULL,
    content     TEXT         NOT NULL,
    is_public   BOOLEAN      NOT NULL,
    top_image   VARCHAR(255) NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
