CREATE TABLE budget_files (
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    budget_id BINARY(16) NOT NULL,
    file_id BINARY(16) NOT NULL,
    UNIQUE uk_budget_id_file_id (budget_id, file_id),
    PRIMARY KEY (id)
);
