CREATE TABLE user_payments
(
    id            BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id       BINARY(16)   NOT NULL,
    `year`        INT          NOT NULL,
    transfer_name VARCHAR(255) NOT NULL,
    checked       BOOLEAN      NOT NULL DEFAULT false,
    note          VARCHAR(255) NOT NULL DEFAULT '',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_user_id_year (user_id, `year`)
);
