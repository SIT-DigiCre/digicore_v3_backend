SELECT 
    BIN_TO_UUID(works.id) AS work_id,
    works.name AS work_name,
    works.updated_at AS work_updated_at,
    BIN_TO_UUID(work_users.user_id) AS author_user_id,
    user_profiles.username AS author_username,
    user_profiles.icon_url AS author_icon_url,
    BIN_TO_UUID(work_work_tags.tag_id) AS tag_id,
    work_tags.name AS tag_name
FROM works
/* IF authorId */
INNER JOIN work_users AS filter_users ON filter_users.work_id = works.id 
    AND filter_users.user_id = UUID_TO_BIN(/*authorId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
/* END */
LEFT JOIN work_users ON work_users.work_id = works.id
LEFT JOIN user_profiles ON user_profiles.user_id = work_users.user_id
LEFT JOIN work_work_tags ON work_work_tags.work_id = works.id
LEFT JOIN work_tags ON work_tags.id = work_work_tags.tag_id
ORDER BY works.updated_at DESC
LIMIT 10
/* IF offset */
OFFSET /*offset*/0
/* END */;
