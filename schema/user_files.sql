CREATE TABLE user_files
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id    BINARY(16)   NOT NULL,
    name       VARCHAR(255) NOT NULL,
    k_size     INT          NOT NULL,
    md5_hash   VARCHAR(255) NOT NULL,
    extension  VARCHAR(255) NOT NULL,
    is_public  BOOLEAN      NOT NULL DEFAULT false,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_id_md5_hash_is_public (user_id, md5_hash, is_public),
    PRIMARY KEY (id)
);
