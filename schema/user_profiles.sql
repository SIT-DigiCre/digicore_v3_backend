CREATE TABLE user_profiles
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
    username                VARCHAR(255) NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(255) NOT NULL DEFAULT '',
    short_introduction VARCHAR(255) NOT NULL DEFAULT 'デジクリ入りました',
    introduction       TEXT         NOT NULL,
    PRIMARY KEY (id)
);
