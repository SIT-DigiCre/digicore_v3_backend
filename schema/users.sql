CREATE TABLE users
(
    id                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    student_number          VARCHAR(8)   NOT NULL,
    active_limit            DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE uk_student_number (student_number),
    PRIMARY KEY (id)
);
