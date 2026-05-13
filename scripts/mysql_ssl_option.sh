#!/usr/bin/env bash
# mysql クライアントの種類に応じて、同一 Docker Network 内 DB 接続用の TLS 無効化オプションを出力する。

set -e

MYSQL_HELP="$(mysql --help 2>/dev/null || true)"

if grep -q -- "--ssl-mode" <<< "${MYSQL_HELP}"; then
  echo "--ssl-mode=DISABLED"
elif grep -q -- "--skip-ssl" <<< "${MYSQL_HELP}"; then
  echo "--skip-ssl"
else
  echo "--ssl=0"
fi
