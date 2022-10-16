-- +migrate Up
ALTER TABLE user_profiles
    CHANGE COLUMN short_self_introduction `short_introduction` VARCHAR(255) NOT NULL DEFAULT 'デジクリ入りました';
ALTER TABLE user_profiles
    CHANGE COLUMN self_introduction `introduction` TEXT NOT NULL;

-- +migrate Down
