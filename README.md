# digicore v3 backend

## 環境構築

1. `.env.sample` をコピーして `.env` を作成する
1. [Discord developers](https://discord.com/developers/applications)で App を作成し、Oauth2 の Redirects に`${FRONTEND_ROOT_URL}/user/discord/callback`を指定する
1. 上記で作成した App の Client information から Client ID と Client Secret を取得し、.env に追記する。
1. [Google Cloud Platform](https://console.cloud.google.com/home/dashboard)で App を作成し、OAuth クライアント ID をアプリケーションの種類をウェブアプリケーションにして作成し、承認済みのリダイレクト URI に`${FRONTEND_ROOT_URL}/signup/callback`と`${FRONTEND_ROOT_URL}/login/callback`を指定する。
1. 上記で作成した App の `client_secret_*.json` をダウンロードし、`config/gcp_secret.json` に名前を書き換えこのファイルが有る階層に配置する。
1. [コンテナのビルド](#コンテナのビルド)を行う
1. [実行](#実行)を行う
1. [DB マイグレーション](#DBマイグレーション)を行う

## コンテナのビルド

```sh
make build
# ubuntuの場合はsudo権限が必要
```

## 実行

```sh
make up-d # デタッチモードで起動
```

## DB マイグレーション

```sh
make migrate-dry # dryrun
make migrate
```

## 開発手順

### テストデータの投入

```sh
make insert_test
```

### api パッケージの更新

> [!WARNING]
>
> `./document/_.gen.yml` と `./pkg/api/_.gen.go` は自動生成であるため直接編集しない

```sh
make 
```

## 開発時の JWT 検証の無効化

`.env` の AUTH を disable に書き換えてください。
