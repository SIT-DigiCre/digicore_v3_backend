CREATE TABLE `groups`
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name        VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    joinable    BOOLEAN      NOT NULL DEFAULT false,
    user_count  INT          NOT NULL DEFAULT 0,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
