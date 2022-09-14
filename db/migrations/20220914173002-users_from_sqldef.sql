-- +migrate Up
CREATE TABLE `users`
(
    `id`                      BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    `student_number`          VARCHAR(8)   NOT NULL,
    UNIQUE `uk_student_number` (`student_number`),
    PRIMARY KEY (`id`)
);

-- +migrate Down
