CREATE TABLE activities
(
    id                     BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                BINARY(16)   NOT NULL,
    place                  VARCHAR(255) NOT NULL,
    note                   TEXT         NULL,
    initial_checked_in_at  DATETIME     NOT NULL,
    initial_checked_out_at DATETIME     NULL,
    checked_in_at          DATETIME     NOT NULL,
    checked_out_at         DATETIME     NULL,
    created_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id_place_checked_in_at (user_id, place, checked_in_at DESC)
);
