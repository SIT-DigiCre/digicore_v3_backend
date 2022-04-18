CREATE TABLE UserPayment
(
    id            BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id       BINARY(16)   NOT NULL,
    `year`        INT          NOT NULL,
    transfer_name VARCHAR(255) NOT NULL,
    checked       BOOLEAN      NOT NULL DEFAULT false,
    PRIMARY KEY (id),
    UNIQUE KEY uq_user_id_year (user_id, `year`),
    CONSTRAINT fk_user_id_user_payment FOREIGN KEY (user_id) REFERENCES User(id)
);
