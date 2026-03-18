SELECT
    upl.link_url,
    upl.created_at,
    upl.updated_at
FROM user_profile_links as upl
WHERE upl.user_id = /* user_id */1
ORDER BY upl.created_at DESC
