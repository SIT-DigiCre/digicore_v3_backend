CREATE TABLE MembershipFee
(
    id            VARCHAR(36)  NOT NULL,
    kind          INT          NOT NULL,
    user_id       BINARY(16)   NOT NULL,
    `year`        INT          NOT NULL,
    transfer_name VARCHAR(255) NOT NULL,
    checked       BOOLEAN      NOT NULL DEFAULT false,
    PRIMARY KEY (id),
    UNIQUE (user_id)
);
