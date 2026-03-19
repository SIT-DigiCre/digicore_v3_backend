SELECT COUNT(*) > 0 AS has_claim
FROM groups_users
INNER JOIN group_claims ON groups_users.group_id = group_claims.group_id
WHERE groups_users.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
  AND group_claims.claim = /*claim*/'infra';
