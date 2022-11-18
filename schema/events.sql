CREATE TABLE events
(
    id          BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name        VARCHAR(255) NOT NULL,
    calendar_view BOOLEAN NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (id)
);
