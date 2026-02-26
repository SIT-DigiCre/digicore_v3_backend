CREAT TABLE user_profile_links
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id     BINARY(16)   NOT NULL UNIQUE,
    `url`         VARCHAR(8000) NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id)
)