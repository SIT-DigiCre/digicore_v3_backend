CREATE TABLE User
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    student_number          VARCHAR(8)   NOT NULL,
    UNIQUE (student_number),
    PRIMARY KEY (id)
);
