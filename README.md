# digicore v3 backend

## 環境構築

1. Windowsにて開発をする方は、Ubuntuで開発を行うために[WSLおよびUbuntuの環境構築](#WSLおよびUbuntuの環境構築)をする　
1. [Docker Desktop](https://www.docker.com)をダウンロード・インストールする
1. VSCodeの拡張機能のところから`ms-vscode-remote.remote-wsl`と検索してVSCodeにWSLの拡張機能をインストールする
1. Ubuntuに開発用の[ディレクトリを作成](#ディレクトリを作成)する 
1. `https://github.com/SIT-DigiCre/digicore_v3_frontend.git` と `https://github.com/SIT-DigiCre/digicore_v3_backend.git` を[クローン](#クローン)する
1. クローンしたディレクトリを[VSCodeで編集](#VSCodeで編集)していく
1. `.env.sample` をコピーして `.env` を作成する
1.  `.env`に環境変数を入力する ※環境変数はsysdevの既存メンバーから教えてもらってください
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
git clone https://github.com/SIT-DigiCre/digicore_v3_frontend.git
git clone https://github.com/SIT-DigiCre/digicore_v3_backend.git
```
## VSCodeで編集
```sh
cd digicore_v3_backend
code . #VSCodeを開く
```

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

### テストの実行

Docker 上の DB を使って backend コンテナ内で全テストを実行できます。

```sh
make test
```

DB コンテナが起動していない場合は、先に起動してください。

```sh
docker compose up -d db
make test
```

### テストデータの投入

```sh
make insert_test
```

### 管理者 claim の付与

指定した学籍番号のアカウントに infra claim と account claim を付与する（`make insert_test` 実行後、infra/account グループが存在することが前提）:

```sh
make grant_admin_claims STUDENT_NUMBER=aa230001
```

### api パッケージの更新

> [!WARNING]
>
> `./document/_.gen.yml` と `./pkg/api/_.gen.go` は自動生成であるため直接編集しない

```sh
make generate_api
```

## JWT 署名鍵の設定

JWT の署名・検証には RSA 秘密鍵（`JWT_PRIVATE_KEY`）を使用します。

- **設定しない場合**: 起動のたびにランダムな鍵が生成されるため、開発環境では問題ありませんが、**コンテナを再起動すると既存の全トークンが失効**します。
- **設定する場合**: 以下の手順で鍵を生成し、`.env` に設定してください。本番環境や複数コンテナ構成では設定が必要です。

### 鍵の生成と .env への書き込み

Makefile の `gen_jwt_key` タスクで、鍵を生成して `.env` に自動書き込みできます:

```sh
make gen_jwt_key
```

- 既存の `JWT_PRIVATE_KEY` が設定されている場合は上書きしません（上書きするには `make gen_jwt_key FORCE=1`）
- 手動で生成する場合は、以下のコマンドで PKCS#8 形式の RSA 2048bit 秘密鍵を生成し、PEM の内容を 1 行（改行は `\n` エスケープ）にして `.env` に設定してください:

```sh
# 鍵生成
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out jwt_private.pem

# 1行に変換してコピー
awk 'BEGIN{ORS="\\n"} {sub(/\n$/,"")} 1' jwt_private.pem
```

```env
JWT_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----\nMIIEv...\n-----END PRIVATE KEY-----"
```

> [!WARNING]
>
> - 鍵は **PKCS#8 形式**（`BEGIN PRIVATE KEY`）である必要があります。`openssl genrsa` で生成した PKCS#1 形式（`BEGIN RSA PRIVATE KEY`）は `openssl pkcs8 -topk8 -nocrypt -in pkcs1.pem -out pkcs8.pem` で変換してください。
> - 鍵を変更するとそれまでの全 JWT が無効になるため、デプロイ時は全ユーザーの再ログインが必要です。
> - 秘密鍵は `.env`（git 管理外）にのみ保管してください。

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
