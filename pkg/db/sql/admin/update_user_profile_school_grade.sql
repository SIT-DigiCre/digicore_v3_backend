UPDATE user_profiles
SET school_grade = /*schoolGrade*/1
WHERE user_id = UUID_TO_BIN(/*userId*/'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee');
