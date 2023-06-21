CREATE TABLE event_reservation_users
(
    id             BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    reservation_id BINARY(16)   NOT NULL,
    user_id        BINARY(16)   NOT NULL,
    `comment`      VARCHAR(255) NOT NULL,
    url            VARCHAR(255),
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE uk_reservation_id_user_id (reservation_id, user_id),
    PRIMARY KEY (id)
);
