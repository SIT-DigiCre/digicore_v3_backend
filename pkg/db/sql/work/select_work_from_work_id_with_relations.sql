SELECT 
    BIN_TO_UUID(works.id) AS work_id,
    works.name AS work_name,
    works.description AS work_description,
    BIN_TO_UUID(work_users.user_id) AS author_user_id,
    user_profiles.username AS author_username,
    user_profiles.icon_url AS author_icon_url,
    BIN_TO_UUID(work_work_tags.tag_id) AS tag_id,
    work_tags.name AS tag_name,
    BIN_TO_UUID(work_files.file_id) AS file_id,
    user_files.name AS file_name
FROM works
LEFT JOIN work_users ON work_users.work_id = works.id
LEFT JOIN user_profiles ON user_profiles.user_id = work_users.user_id
LEFT JOIN work_work_tags ON work_work_tags.work_id = works.id
LEFT JOIN work_tags ON work_tags.id = work_work_tags.tag_id
LEFT JOIN work_files ON work_files.work_id = works.id
LEFT JOIN user_files ON user_files.id = work_files.file_id
WHERE works.id = UUID_TO_BIN(/*workId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
