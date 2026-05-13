#!/usr/bin/env bash
# mysql クライアントの種類に応じて、同一 Docker Network 内 DB 接続用の TLS 無効化オプションを出力する。

set -e

if mysql --help 2>/dev/null | grep -q -- "--ssl-mode"; then
  echo "--ssl-mode=DISABLED"
elif mysql --help 2>/dev/null | grep -q -- "--skip-ssl"; then
  echo "--skip-ssl"
else
  echo "--ssl=0"
fi
