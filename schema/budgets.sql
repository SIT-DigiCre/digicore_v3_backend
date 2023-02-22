CREATE TABLE budgets (
    id BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    name VARCHAR(255) NOT NULL,
    proposer_user_id BINARY(16) NOT NULL,
    approver_user_id BINARY(16),
    `status` VARCHAR(255) NOT NULL,
    class VARCHAR(255) NOT NULL,
    budget INT NOT NULL,
    settlement INT NOT NULL,
    purpose TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved_at DATETIME,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    remarks TEXT,
    mattermost_url VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);
