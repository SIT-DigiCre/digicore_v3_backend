CREATE TABLE UserPrivateProfile
(
    id                      VARCHAR(36)  NOT NULL,
    user_id                 VARCHAR(36)  NOT NULL,
    first_name              VARCHAR(255) NOT NULL,
    last_name               VARCHAR(255) NOT NULL,
    first_name_kana         VARCHAR(255) NOT NULL,
    last_name_kana          VARCHAR(255) NOT NULL,
    phone_number            VARCHAR(32)  NOT NULL,
    address                 VARCHAR(255) NOT NULL,
    parent_name             VARCHAR(255) NOT NULL,
    parent_cellphone_number VARCHAR(32)  NOT NULL,
    parent_homephone_number VARCHAR(32)  NOT NULL,
    parent_address          VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (user_id),
    CONSTRAINT fk_user_id_user_private_profile FOREIGN KEY (user_id) REFERENCES User(id)
);
