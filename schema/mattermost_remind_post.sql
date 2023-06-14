CREATE TABLE mattermost_remind_post
(
    id           BINARY(16)  NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_name    VARCHAR(255) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    body         TEXT        NOT NULL,
    remind_date  DATETIME    NOT NULL,
    posted       BOOLEAN     NOT NULL DEFAULT false,
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
