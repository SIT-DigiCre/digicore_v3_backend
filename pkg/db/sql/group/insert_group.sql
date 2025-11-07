INSERT INTO `groups` (id, name, description, joinable, user_count, created_at, updated_at)
VALUES (UUID_TO_BIN(@id), /*name*/'', /*description*/'', /*joinable*/true, 1, NOW(), NOW());
