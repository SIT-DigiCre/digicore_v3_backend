CREATE TABLE UserProfile
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
    username                VARCHAR(255) NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(255) NOT NULL DEFAULT '',
    active_limit            DATE         NOT NULL DEFAULT (CURRENT_DATE),
    short_self_introduction VARCHAR(255) NOT NULL,
    self_introduction       VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT fk_user_id_user_profile FOREIGN KEY (user_id) REFERENCES User (id)
);
