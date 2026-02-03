-- 管理者用ユーザー一覧の総件数取得クエリ
-- select_admin_user_list.sqlと同じフィルタ条件を適用し、ページネーション用の総件数を返す
-- サブクエリでGROUP BY後の件数をカウントする必要がある（HAVINGを使うため）

SELECT COUNT(*) as total FROM (
    SELECT
        user_profiles.user_id
    FROM user_profiles
    LEFT JOIN users ON users.id = user_profiles.user_id
    -- 管理者判定用
    LEFT JOIN groups_users gu ON gu.user_id = user_profiles.user_id
    LEFT JOIN group_claims gc ON gc.group_id = gu.group_id AND gc.claim = 'admin'
    WHERE 1 = 1
    /* IF query */
      AND (
        username LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
        OR users.student_number LIKE CONCAT('%', REPLACE(REPLACE(/*query*/'', '%', '\%'), '_', '\_'), '%')
      )
    /* END */
    /* IF schoolGrade */
      AND school_grade = /*schoolGrade*/0
    /* END */
    GROUP BY user_profiles.user_id
    /* IF isAdmin */
      HAVING MAX(CASE WHEN gc.claim IS NOT NULL THEN true ELSE false END) = /*isAdmin*/false
    /* END */
) cnt_sub;
