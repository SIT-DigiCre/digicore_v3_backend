CREATE TABLE event_reservation_users
(
    id             BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    reservation_id BINARY(16)   NOT NULL,
    user_id        BINARY(16)   NOT NULL,
    PRIMARY KEY (id)
);
