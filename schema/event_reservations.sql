CREATE TABLE event_reservations
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    event_id                BINARY(16)   NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    description             TEXT         NOT NULL,
    start_date              DATETIME     NOT NULL,
    finish_date             DATETIME     NOT NULL,
    reservation_start_date  DATETIME     NOT NULL,
    reservation_finish_date DATETIME     NOT NULL,
    capacity                INT          NOT NULL,
    PRIMARY KEY (id)
);
