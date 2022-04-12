CREATE TABLE User
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    student_number          VARCHAR(8)   NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE `Group`
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name        VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE GroupUser
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id  BINARY(16) NOT NULL,
    group_id BINARY(16) NOT NULL,
    PRIMARY KEY (id)
);

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

CREATE TABLE UserPrivateProfile
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
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
    CONSTRAINT fk_user_id_user_private_profile FOREIGN KEY (user_id) REFERENCES User(id)
);

CREATE TABLE UserProfile
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                 BINARY(16)   NOT NULL UNIQUE,
    username                VARCHAR(255) NOT NULL,
    school_grade            INT          NOT NULL,
    icon_url                VARCHAR(255) NOT NULL,
    discord_userid          VARCHAR(255) NOT NULL DEFAULT '',
    active_limit            DATE         NOT NULL DEFAULT (CURRENT_DATE),
    short_self_introduction VARCHAR(255) NOT NULL,
    self_introduction       VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT fk_user_id_user_profile FOREIGN KEY (user_id) REFERENCES User (id)
);
