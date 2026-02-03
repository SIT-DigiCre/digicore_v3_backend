# AGENTS.md

このファイルは、AI Agents がこのリポジトリで作業する際の指針を提供します。

**重要**: AI Agents とのやり取りは常に日本語で行ってください。コードコメントやコミットメッセージも日本語で書いてください。

## プロジェクト概要

Digicore v3 Backend は、日本のデジタルクリエイション学生団体の管理システム用 Go REST API です。ユーザー認証、イベント管理、予算追跡、作品ポートフォリオ、各種統合機能を提供します。

## 技術スタック

- **言語**: Go 1.19+
- **フレームワーク**: Echo v4 web framework
- **データベース**: MySQL 8.0.33
- **認証**: JWT with Google/Discord OAuth
- **API**: OpenAPI 3.0.3 code generation
- **SQL**: go-twowaysql templating engine
- **コンテナ**: Docker with hot-reload for development

## よく使うコマンド

### 開発環境

```bash
# 開発環境のビルドと起動
make build && make up

# 環境停止
make down

# ログ確認
make logs

# ホットリロード開発（backendコンテナ内で実行）
air -c ./config/.air.toml
```

### API コード生成

```bash
# OpenAPI仕様からAPIコードを生成（document/*.ymlを編集後に必須）
make generate_api
```

### データベース操作

```bash
# マイグレーションテスト（dry run）
make migrate-dry

# マイグレーション実行
make migrate

# テストデータ挿入
make insert_test
```

### 認証設定

- 本番: JWT 認証有効
- 開発: テスト用に`.env`で`AUTH=disable`に設定

### テスト用 JWT トークン生成

`AUTH=disable` でも API リクエストには JWT トークンが必要です（user_id の取得に使用）。

```bash
# テスト用トークン生成スクリプト
./scripts/gen_jwt.sh              # デフォルト user_id
./scripts/gen_jwt.sh <user_id>    # 指定 user_id

# API テスト例
TOKEN=$(./scripts/gen_jwt.sh)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8000/user/me
```

## アーキテクチャ

### コード構成

- **ドメイン駆動設計**: ビジネスドメイン別に構成（`pkg/user/`, `pkg/event/`, `pkg/budget/`など）
- **クリーンアーキテクチャ**: API 層 → ビジネスロジック → データアクセス
- **コード生成**: OpenAPI 仕様による型安全な API コード生成

### 主要パターン

**HTTP ハンドラ**: 各ドメインで統一された命名規則:

- `get_*.go`: 読み取り操作
- `post_*.go`: 作成操作
- `put_*.go`: 更新操作
- `delete_*.go`: 削除操作

**データベースアクセス**:

- 全 SQL クエリは`pkg/db/sql/`にドメイン別で管理
- `twowaysql`テンプレートエンジンと struct タグでパラメータ指定
- データベースインターフェース: `db.Client`の`Select()`と`Exec()`メソッド
- `TransactionClient`によるトランザクション支援

**エラーハンドリング**: `pkg/api/response`による HTTP ステータス、ログレベル、日本語ユーザーメッセージの構造化エラー

### コード生成ワークフロー

1. `document/`ディレクトリの OpenAPI 仕様を編集
2. `make generate_api`で生成されたコードを更新
3. **重要**: `*.gen.go`や`*.gen.yml`ファイルは直接編集禁止

### 環境セットアップ

1. `.env.sample`を`.env`にコピー
2. Discord/Google OAuth 認証情報を設定
3. Google OAuth JSON を`config/gcp_secret.json`に配置
4. `make build && make up`を実行
5. データベースセットアップのため`make migrate`を実行

## ドメイン構造

各ドメインパッケージ（user, event, budget, work, group など）の構成:

- ビジネスロジック関数
- OpenAPI エンドポイントに対応する HTTP ハンドラ実装
- `pkg/db/sql/{domain}/`の専用 SQL クエリファイル

## 外部統合

- Google OAuth（ログイン/サインアップ）
- Discord OAuth と Webhook
- Mattermost チーム管理
- AWS S3/Wasabi オブジェクトストレージ
