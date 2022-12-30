CREATE TABLE blog_blog_tags
(
    id         BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    blog_id    BINARY(16)   NOT NULL,
    tag_id     BINARY(16)   NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_blog_id_tag_id (blog_id, tag_id),
    PRIMARY KEY (id)
);
