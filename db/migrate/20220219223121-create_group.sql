-- +migrate Up
CREATE TABLE `Group`
(
    id          BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name        VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    PRIMARY KEY (id)
);


-- +migrate Down
DROP TABLE `Group`;
