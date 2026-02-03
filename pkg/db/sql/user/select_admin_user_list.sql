-- 管理者用ユーザー一覧取得クエリ
-- ユーザーの公開プロフィール、非公開プロフィール、管理者フラグを取得する
--
-- テーブル結合:
--   user_profiles: ユーザーの公開プロフィール（必須）
--   users: 学籍番号（LEFT JOIN: プロフィールのみ存在するケースに対応）
--   user_private_profiles: 非公開プロフィール（LEFT JOIN: 未登録の場合はNULL→COALESCEで空文字に変換）
--   groups_users + group_claims: 管理者権限の判定用（adminクレームを持つグループに所属しているか）
--
-- フィルタ条件（twowaysql動的条件）:
--   query: ユーザー名または学籍番号の部分一致検索（%と_はエスケープ）
--   schoolGrade: 学年での絞り込み
--   isAdmin: 管理者フラグでの絞り込み（HAVING句で集計後にフィルタ）

SELECT
    BIN_TO_UUID(user_profiles.user_id) as user_id,
    users.student_number,
    username,
    school_grade,
    icon_url,
    discord_userid,
    active_limit,
    short_introduction,
    introduction,
    -- adminクレームを持つグループに1つでも所属していれば管理者
    MAX(CASE WHEN gc.claim IS NOT NULL THEN true ELSE false END) as is_admin,
    -- 非公開プロフィールはNULLの可能性があるためCOALESCEでデフォルト値を設定
    COALESCE(first_name, '') as first_name,
    COALESCE(last_name, '') as last_name,
    COALESCE(first_name_kana, '') as first_name_kana,
    COALESCE(last_name_kana, '') as last_name_kana,
    COALESCE(is_male, false) as is_male,
    COALESCE(phone_number, '') as phone_number,
    COALESCE(address, '') as address,
    COALESCE(parent_name, '') as parent_name,
    COALESCE(parent_last_name, '') as parent_last_name,
    COALESCE(parent_first_name, '') as parent_first_name,
    COALESCE(parent_cellphone_number, '') as parent_cellphone_number,
    COALESCE(parent_homephone_number, '') as parent_homephone_number,
    COALESCE(parent_address, '') as parent_address
FROM user_profiles
LEFT JOIN users ON users.id = user_profiles.user_id
LEFT JOIN user_private_profiles ON user_private_profiles.user_id = user_profiles.user_id
-- 管理者判定用: ユーザーが所属するグループのクレームを取得
LEFT JOIN groups_users gu ON gu.user_id = user_profiles.user_id
LEFT JOIN group_claims gc ON gc.group_id = gu.group_id AND gc.claim = 'admin'
WHERE 1 = 1
/* IF query */
  -- ユーザー名または学籍番号で部分一致検索（ワイルドカード文字をエスケープ）
  AND (
    username LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
    OR users.student_number LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
  )
/* END */
/* IF schoolGrade */
  AND school_grade = /*schoolGrade*/0
/* END */
-- 複数グループに所属する場合があるためGROUP BYで集約
GROUP BY user_profiles.user_id
/* IF isAdmin */
  -- 集計後に管理者フラグでフィルタ
  HAVING MAX(CASE WHEN gc.claim IS NOT NULL THEN true ELSE false END) = /*isAdmin*/false
/* END */
ORDER BY user_profiles.created_at DESC
LIMIT /*limit*/100
OFFSET /*offset*/0;
