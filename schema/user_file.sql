CREATE TABLE UserFile
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id    BINARY(16)   NOT NULL,
    name       VARCHAR(255) NOT NULL,
    k_size     INT          NOT NULL,
    md5_hash   VARCHAR(255) NOT NULL,
    extension  VARCHAR(255) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_user_id_user_file FOREIGN KEY (user_id) REFERENCES User(id)
);
