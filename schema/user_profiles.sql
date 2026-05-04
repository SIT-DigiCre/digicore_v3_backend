CREATE TABLE user_profiles
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
    username                VARCHAR(255) NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(255) NOT NULL DEFAULT '',
    active_limit            DATE         NOT NULL,
    is_graduated            BOOLEAN      NOT NULL DEFAULT false,
    is_member               BOOLEAN      NOT NULL DEFAULT true,
    short_introduction VARCHAR(255) NOT NULL DEFAULT 'デジクリ入りました',
    introduction       TEXT         NOT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
