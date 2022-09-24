# digicore v3 backend

## 環境構築

1. .env.sampleをコピーして.envを作成する
1. [Discord developers](https://discord.com/developers/applications)でAppを作成し、Oauth2のRedirectsに`${BACKEND_ROOT_URL}/discord/oauth/callback`を指定する
1. 上記で作成したAppのClient informationからClient IDとClient Secretを取得し、.envに追記する。
1. [Google Cloud Platform](https://console.cloud.google.com/home/dashboard)でAppを作成し、OAuth クライアント IDをアプリケーションの種類をウェブアプリケーションにして作成し、承認済みのリダイレクト URIに`${FRONTEND_ROOT_URL}/signup/callback`と`${FRONTEND_ROOT_URL}/login/callback`を指定する。
1. 上記で作成したAppのclient_secret_*.jsonをダウンロードし、client_secret.jsonに名前を書き換えこのファイルが有る階層に配置する。
1. [コンテナのビルド](#コンテナのビルド)を行う
1. [実行](#実行)を行う
1. [DBマイグレーション](#DBマイグレーション)を行う

## コンテナのビルド

```sh
docker compose build
```

## 実行

```sh
docker compose up
```

## DBマイグレーション

```bash
docker compose run --rm -w /app/db admin sql-migrate up
```

## 開発手順

### apiパッケージの更新

**./document/bundle.ymlと./pkg/api/*.gen.goは自動生成であるため直接編集しない**

```sh
docker compose run --rm -w /app/document node_tool swagger-cli bundle -o ./bundle.yml -t yaml ./openapi.yml # OpenAPIファイルの結合
docker compose run --rm -w /app admin make generate_api # apiパッケージの生成
```

## 開発時のJWT検証の無効化

.envのAUTHをdisableに書き換えてください。
