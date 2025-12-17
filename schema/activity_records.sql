CREATE TABLE activity_records
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id    BINARY(16)   NOT NULL,
    place      VARCHAR(255) NOT NULL,
    type       VARCHAR(255) NOT NULL,
    datetime   DATETIME     NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id_place_datetime (user_id, place, datetime DESC)
);


