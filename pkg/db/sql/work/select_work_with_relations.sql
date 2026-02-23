SELECT
    BIN_TO_UUID(works.id) AS work_id,
    works.name AS work_name,
    works.description AS work_description,
    BIN_TO_UUID(work_users.user_id) AS author_user_id,
    user_profiles.username AS author_username,
    user_profiles.icon_url AS author_icon_url,
    BIN_TO_UUID(work_work_tags.tag_id) AS tag_id,
    work_tags.name AS tag_name,
    BIN_TO_UUID(first_work_file.file_id) AS file_id,
    user_files.name AS file_name
FROM (
    SELECT works.id, works.name, works.description, works.updated_at
    FROM works
    /* IF authorId */
    INNER JOIN work_users AS filter_users ON filter_users.work_id = works.id 
        AND filter_users.user_id = UUID_TO_BIN(/*authorId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
    /* END */
    ORDER BY works.updated_at DESC
    LIMIT 10
    /* IF offset */
    OFFSET /*offset*/0
    /* END */
) AS works
LEFT JOIN work_users ON work_users.work_id = works.id
LEFT JOIN user_profiles ON user_profiles.user_id = work_users.user_id
LEFT JOIN work_work_tags ON work_work_tags.work_id = works.id
LEFT JOIN work_tags ON work_tags.id = work_work_tags.tag_id
LEFT JOIN (
    SELECT ranked.work_id, ranked.file_id
    FROM (
        SELECT
            work_files.work_id,
            work_files.file_id,
            ROW_NUMBER() OVER (
                PARTITION BY work_files.work_id
                ORDER BY work_files.created_at ASC, work_files.id ASC
            ) AS row_num
        FROM work_files
    ) AS ranked
    WHERE ranked.row_num = 1
) AS first_work_file ON first_work_file.work_id = works.id
LEFT JOIN user_files ON user_files.id = first_work_file.file_id
ORDER BY works.updated_at DESC;
