SELECT
    BIN_TO_UUID(user_profiles.user_id) as user_id,
    student_number,
    username,
    school_grade,
    icon_url,
    discord_userid,
    active_limit,
    is_graduated,
    is_member,
    short_introduction
FROM user_profiles
LEFT JOIN users ON users.id = user_profiles.user_id
WHERE user_profiles.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
