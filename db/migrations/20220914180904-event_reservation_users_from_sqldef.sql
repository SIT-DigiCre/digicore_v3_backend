-- +migrate Up
CREATE TABLE event_reservation_users
(
    `id`             BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    `reservation_id` BINARY(16)   NOT NULL,
    `user_id`        BINARY(16)   NOT NULL,
    `comment`        VARCHAR(255) NOT NULL,
    `url`            VARCHAR(255),
    PRIMARY KEY (`id`)
);

-- +migrate Down
