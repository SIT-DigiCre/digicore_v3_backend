CREATE TABLE UserPrivateProfile
(
    id                      VARCHAR(36)  NOT NULL,
    user_id                 VARCHAR(36)  NOT NULL,
    first_name              VARCHAR(255) NOT NULL DEFAULT '',
    last_name               VARCHAR(255) NOT NULL DEFAULT '',
    first_name_kana         VARCHAR(255) NOT NULL DEFAULT '',
    last_name_kana          VARCHAR(255) NOT NULL DEFAULT '',
    phone_number            VARCHAR(32)  NOT NULL DEFAULT '',
    address                 VARCHAR(255) NOT NULL DEFAULT '',
    parent_name             VARCHAR(255) NOT NULL DEFAULT '',
    parent_cellphone_number VARCHAR(32)  NOT NULL DEFAULT '',
    parent_homephone_number VARCHAR(32)  NOT NULL DEFAULT '',
    parent_address          VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    UNIQUE (user_id),
    CONSTRAINT fk_user_id_user_private_profile FOREIGN KEY (user_id) REFERENCES User(id)
);
