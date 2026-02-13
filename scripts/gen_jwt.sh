#!/usr/bin/env bash
# AUTH=disable 時のテスト用 JWT を生成する。
# 署名は検証されないが、base64url として有効な文字列にする必要がある（.dummy は長さで失敗するため）。
#
# 使い方:
#   ./scripts/gen_jwt.sh                    # デフォルト user_id (11111111-...)
#   ./scripts/gen_jwt.sh <user_id>          # 指定 user_id
#   ./scripts/gen_jwt.sh <user_id> <exp>   # exp も指定（Unix 秒）

set -e

USER_ID="${1:-11111111-1111-1111-1111-111111111111}"
EXP="${2:-9999999999}"

base64url() {
  base64 | tr -d '\n' | tr '+/' '-_' | tr -d '='
}

HEADER=$(echo -n '{"alg":"RS256","typ":"JWT"}' | base64url)
PAYLOAD=$(echo -n "{\"sub\":\"$USER_ID\",\"exp\":$EXP}" | base64url)
# 署名は AUTH=disable 時は検証されないが、デコードできる有効な base64url にすること
SIG="c2ln"
echo "${HEADER}.${PAYLOAD}.${SIG}"
