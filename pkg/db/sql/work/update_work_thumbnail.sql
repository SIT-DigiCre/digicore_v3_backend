UPDATE works SET thumbnail_id = UUID_TO_BIN(NULLIF(/*thumbnailId*/'', ''))
WHERE id = UUID_TO_BIN(/*workId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
