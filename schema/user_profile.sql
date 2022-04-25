CREATE TABLE UserProfile
(
    id                      VARCHAR(36)  NOT NULL,
    user_id                 VARCHAR(36)  NOT NULL,
    username                VARCHAR(255) NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(255) NOT NULL DEFAULT '',
    active_limit            DATE         NOT NULL,
    short_self_introduction VARCHAR(255) NOT NULL DEFAULT 'デジクリ入りました',
    self_introduction       VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    UNIQUE (user_id),
    CONSTRAINT fk_user_id_user_profile FOREIGN KEY (user_id) REFERENCES User (id)
);
