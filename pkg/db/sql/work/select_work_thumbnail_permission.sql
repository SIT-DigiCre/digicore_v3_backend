SELECT (t.user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee')
    OR EXISTS (SELECT 1 FROM works w WHERE w.id = UUID_TO_BIN(/*workId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee') AND w.thumbnail_id = t.id)) AS allowed
FROM work_thumbnails t
WHERE t.id = UUID_TO_BIN(/*thumbnailId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
