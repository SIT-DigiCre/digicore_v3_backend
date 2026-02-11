CREATE TABLE user_private_profiles
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
    first_name              VARCHAR(255) NOT NULL,
    last_name               VARCHAR(255) NOT NULL,
    first_name_kana         VARCHAR(255) NOT NULL,
    last_name_kana          VARCHAR(255) NOT NULL,
    is_male                 BOOLEAN      NOT NULL DEFAULT true,
    phone_number            VARCHAR(32)  NOT NULL,
    address                 VARCHAR(255) NOT NULL,
    parent_name             VARCHAR(255) NOT NULL,
    parent_last_name        VARCHAR(255) NOT NULL DEFAULT '',
    parent_first_name       VARCHAR(255) NOT NULL DEFAULT '',
    parent_cellphone_number VARCHAR(32)  NOT NULL,
    parent_homephone_number VARCHAR(32)  NOT NULL DEFAULT '',
    parent_address          VARCHAR(255) NOT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
