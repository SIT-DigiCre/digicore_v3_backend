SELECT BIN_TO_UUID(id) AS blog_id, name, content, top_image, is_public, created_at, updated_at FROM blogs WHERE id = UUID_TO_BIN(/*blogId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
