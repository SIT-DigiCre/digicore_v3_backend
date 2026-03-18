DELETE FROM user_profile_links
WHERE link_url = /*link_url*/'https://example.com' AND user_id = UUID_TO_BIN(/*user_id*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
