CREATE TABLE group_claims
(
    id       BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    group_id BINARY(16) NOT NULL,
    claim    VARCHAR(255) NOT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE uk_user_id_group_id (group_id, claim),
    PRIMARY KEY (id)
);
