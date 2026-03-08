#!/usr/bin/env bash
# 指定した student_number のアカウントに infra claim と account claim を付与する。
# groups_users に該当グループへの参加レコードを挿入する。
#
# 前提: make insert_test で infra/account グループが作成済みであること。
# 実行: docker compose 環境で make grant_admin_claims STUDENT_NUMBER=<学籍番号>
#   または: docker compose run --rm admin bash scripts/grant_admin_claims.sh <学籍番号>
#
# 使い方:
#   ./scripts/grant_admin_claims.sh aa230001
#   ./scripts/grant_admin_claims.sh bp230002

set -e

STUDENT_NUMBER="${1:?学籍番号を指定してください。例: ./scripts/grant_admin_claims.sh aa230001}"

# 学籍番号の簡易バリデーション（英数字とハイフンのみ許可）
if ! [[ "$STUDENT_NUMBER" =~ ^[a-zA-Z0-9\-]+$ ]]; then
  echo "エラー: 学籍番号は英数字のみ使用できます: $STUDENT_NUMBER"
  exit 1
fi

# DB接続情報（docker compose 実行時に .env から注入される）
MYSQL_CMD="mysql -u ${DB_USER} -p${DB_PASSWORD} --host=${DB_HOST} ${DB_DATABASE} --default-character-set=utf8mb4 -N"

# 1. student_number から user_id を取得
USER_ID=$($MYSQL_CMD -e "SELECT BIN_TO_UUID(id) FROM users WHERE student_number = '${STUDENT_NUMBER}' LIMIT 1;" 2>/dev/null || true)

if [ -z "$USER_ID" ]; then
  echo "エラー: 学籍番号 '$STUDENT_NUMBER' に該当するユーザーが見つかりません。"
  exit 1
fi

echo "ユーザーを検出: user_id=$USER_ID (student_number=$STUDENT_NUMBER)"

# 2. infra / account の group_id を取得
INFRA_GROUP_ID=$($MYSQL_CMD -e "SELECT BIN_TO_UUID(group_id) FROM group_claims WHERE claim = 'infra' LIMIT 1;" 2>/dev/null || true)
ACCOUNT_GROUP_ID=$($MYSQL_CMD -e "SELECT BIN_TO_UUID(group_id) FROM group_claims WHERE claim = 'account' LIMIT 1;" 2>/dev/null || true)

if [ -z "$INFRA_GROUP_ID" ]; then
  echo "エラー: infra claim を持つグループが存在しません。make insert_test を実行してください。"
  exit 1
fi

if [ -z "$ACCOUNT_GROUP_ID" ]; then
  echo "エラー: account claim を持つグループが存在しません。make insert_test を実行してください。"
  exit 1
fi

# 3. groups_users に挿入（既に存在する場合は IGNORE）
INSERT_SQL="
INSERT IGNORE INTO groups_users (id, group_id, user_id) VALUES
  (UUID_TO_BIN(UUID()), UUID_TO_BIN('${INFRA_GROUP_ID}'), UUID_TO_BIN('${USER_ID}')),
  (UUID_TO_BIN(UUID()), UUID_TO_BIN('${ACCOUNT_GROUP_ID}'), UUID_TO_BIN('${USER_ID}'));
"

$MYSQL_CMD -e "$INSERT_SQL"

echo "infra claim と account claim を付与しました。"
echo "  - infra グループ: $INFRA_GROUP_ID"
echo "  - account グループ: $ACCOUNT_GROUP_ID"
