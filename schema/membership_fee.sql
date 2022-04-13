CREATE TABLE MembershipFee
(
    id            BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    kind          INT          NOT NULL,
    user_id       BINARY(16)   NOT NULL,
    `year`          INT          NOT NULL,
    transfer_name VARCHAR(255) NOT NULL,
    checked       BOOLEAN      NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);
