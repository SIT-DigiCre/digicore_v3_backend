CREATE TABLE `groups`
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name        VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    joinable    BOOLEAN      NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);
