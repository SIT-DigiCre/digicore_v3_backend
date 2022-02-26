-- +migrate Up
CREATE TABLE User
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    username                VARCHAR(255) NOT NULL,
    student_number          VARCHAR(8)   NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(32)  NOT NULL,
    active_limit            DATE         NOT NULL,
    first_name              VARCHAR(255) NOT NULL,
    last_name               VARCHAR(255) NOT NULL,
    first_name_kana         VARCHAR(255) NOT NULL,
    last_name_kana          VARCHAR(255) NOT NULL,
    phone_number            VARCHAR(255) NOT NULL,
    address                 VARCHAR(255) NOT NULL,
    parent_name             VARCHAR(255) NOT NULL,
    parent_cellphone_number VARCHAR(255) NOT NULL,
    parent_homephone_number VARCHAR(255),
    parent_address          VARCHAR(255) NOT NULL,
    self_introduction       VARCHAR(255) NOT NULL,
    short_self_introduction TEXT         NOT NULL,
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE User;
