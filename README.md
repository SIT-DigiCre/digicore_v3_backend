# digicore v3 backend

## 環境構築

1. Windowsにて開発をする方は、Ubuntuで開発を行うために[WSLおよびUbuntuの環境構築](#WSLおよびUbuntuの環境構築)をする　
1. [Docker Desktop](https://www.docker.com)をダウンロード・インストールする
1. VSCodeの拡張機能のところからWSLと検索してVSCodeにWSLの拡張機能をインストールする
1. Ubuntuに開発用の[ディレクトリを作成](#ディレクトリを作成)する 
1. `https://github.com/SIT-DigiCre/digicore_v3_frontend.git` と `https://github.com/SIT-DigiCre/digicore_v3_backend.git` を[クローン](#クローン)する
1. クローンしたディレクトリを[VSCodeで編集](#VSCodeで編集)していく
1. `.env.sample` をコピーして `.env` を作成する
1.  `.env`に環境変数を入力する
1. [Discord developers](https://discord.com/developers/applications)で App を作成し、Oauth2 の Redirects に`${FRONTEND_ROOT_URL}/user/discord/callback`を指定する
1. 上記で作成した App の Client information から Client ID と Client Secret を取得し、.env に追記する。
1. [Google Cloud Platform](https://console.cloud.google.com/home/dashboard)で App を作成し、OAuth クライアント ID をアプリケーションの種類をウェブアプリケーションにして作成し、承認済みのリダイレクト URI に`${FRONTEND_ROOT_URL}/signup/callback`と`${FRONTEND_ROOT_URL}/login/callback`を指定する。
1. 上記で作成した App の `client_secret_*.json` をダウンロードし、`config/gcp_secret.json` に名前を書き換えこのファイルが有る階層に配置する。
1. [コンテナのビルド](#コンテナのビルド)を行う
1. [実行](#実行)を行う
1. [DB マイグレーション](#DBマイグレーション)を行う

##  WSLおよびUbuntuの環境構築
```sh
wsl --install
wsl --install -d Ubuntu-24.04.4 
```
## ディレクトリを作成
```sh
cd    ##ホームディレクトリに戻る
mkdir digicre
cd digicre #digicreフォルダに入る
```

## クローン
```sh
git clone 指定したURL
#今開いているリポジトリに指定したURLのリポジトリを複製する
```
## VSCodeで編集
```sh
cd digicore_v3_backend
code . #VSCodeを開く
```

## コンテナのビルド

```sh
make build
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
make generate_api
```

## 開発時の JWT 検証の無効化

`.env` の AUTH を disable に書き換えてください。

## リクエストのテスト

### 認証ヘッダー設定

認証が無効化されていても、JWT トークンの subject から user_id を取得するため、有効な JWT 形式のトークンが必要です。以下のようなダミートークンを生成：

```sh
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTEiLCJleHAiOjk5OTk5OTk5OTl9.dummy
```

トークンの生成は https://www.jwt.io/ja から行ってください。

#### ヘッダーの例

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

#### ペイロードの例

```json
{
  "sub": "55555555-5555-5555-5555-555555555555"
}
```
