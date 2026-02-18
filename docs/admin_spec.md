# 管理者権限・グループ管理 仕様書

## 管理者権限モデル

### 役職 claim 一覧

| claim | 役職 |
|-------|------|
| `account` | 会計 |
| `infra` | インフラ |

これらは `pkg/admin/admin.go` の `adminClaims` 変数で管理されており、`GetAdminClaims()` でコピーを取得できます。

```go
var adminClaims = []string{"account", "infra"}

func GetAdminClaims() []string { ... }
```

### 「管理者」の定義

`GET /user/me` のレスポンスに含まれる `isAdmin` フラグは、`AdminClaims` のいずれかの claim を持つユーザーで `true` になります（`CheckUserIsAdmin`）。

---

## 認証・認可アーキテクチャ

### OpenAPI スコープによる権限制御

権限制御は OpenAPI の `security` フィールドで管理します。

```yaml
# 全認証ユーザーが利用可能
security:
  - BearerAuth: []

# infra claim が必要
security:
  - BearerAuth:
      - "infra"

# account claim が必要
security:
  - BearerAuth:
      - "account"
```

本番環境では `OapiRequestValidatorWithOptions` ミドルウェアがスコープを検証します。

### コード内での追加チェック

OpenAPI スコープで表現できない条件付き認可が必要な場合のみ、ビジネスロジック層でのチェックを行います。原則として二重チェック（OpenAPI スコープ + コード内チェック）は避けます。

**例外: `PUT /activity/record/{recordId}`**

- レコードの所有者: 誰でも編集可能
- 他ユーザーのレコード: `infra` claim が必要

この条件分岐は OpenAPI スコープでは表現できないため、`pkg/activity/put_activity_record_record_id.go` でコード内チェックを行います。

---

## claim 別エンドポイント一覧

### `infra` claim が必要なエンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| `POST` | `/group/admin` | 任意の claim を持つグループを作成 |
| `POST` | `/activity/checkout/{userId}` | 管理者による強制チェックアウト |
| `PUT` | `/activity/record/{recordId}` | 他ユーザーの在室レコードを編集（所有者本人は不要） |
| `POST` | `/mail` | メール一括送信（アドレス直接指定またはユーザーID指定） |

### `account` claim が必要なエンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| `GET` | `/payment` | 部費支払い一覧取得 |
| `GET` | `/payment/{paymentId}` | 部費支払い詳細取得 |
| `PUT` | `/payment/{paymentId}` | 部費支払い状況の更新 |
| `PUT` | `/budget/{budgetId}/admin` | 稟議の承認・ステータス変更 |

---

## グループ機能

### エンドポイント一覧

| メソッド | パス | 必要な claim | 説明 |
|---------|------|-------------|------|
| `GET` | `/group` | なし（要認証） | グループ一覧取得 |
| `POST` | `/group` | なし（要認証） | claim なし通常グループ作成 |
| `POST` | `/group/admin` | `infra` | 任意の claim を持つグループ作成 |
| `POST` | `/group/{groupId}/user` | なし（要グループ所属 or `infra`） | グループにユーザーを追加 |
| `POST` | `/group/{groupId}/join` | なし（要認証） | グループに自発参加 |

### グループ作成の権限ルール

- **claim なしのグループ** (`POST /group`): 認証済みユーザー全員が作成可能
- **任意の claim を持つグループ** (`POST /group/admin`): `infra` claim を持つユーザーのみ作成可能

### 自発参加の制限

`POST /group/{groupId}/join` で参加できないケース:

- `AdminClaims` のいずれかの claim を持つグループ（`account` / `infra`）→ 403（`joinable` フラグに関わらず無条件で拒否）
- `joinable=false` のグループ → 403

---

## データベース構造

### `groups` テーブル

| カラム | 型 | 説明 |
|--------|----|------|
| `id` | BINARY(16) | グループ ID（UUID） |
| `name` | VARCHAR | グループ名 |
| `description` | TEXT | 説明 |
| `joinable` | BOOL | 自発参加可否 |
| `user_count` | INT | メンバー数（非正規化） |

### `groups_users` テーブル

| カラム | 型 | 説明 |
|--------|----|------|
| `id` | BINARY(16) | レコード ID |
| `group_id` | BINARY(16) | グループ ID（FK） |
| `user_id` | BINARY(16) | ユーザー ID（FK） |

ユニーク制約: `(group_id, user_id)`

### `group_claims` テーブル

| カラム | 型 | 説明 |
|--------|----|------|
| `id` | BINARY(16) | レコード ID |
| `group_id` | BINARY(16) | グループ ID（FK） |
| `claim` | VARCHAR | claim 文字列（例: `"infra"`, `"account"`） |

ユニーク制約: `(group_id, claim)`

---

## テストデータ

`make insert_test` で投入されるシードデータのユーザーと権限状態。

| user_id | 名前 | 学籍番号 | 役職 claim |
|---------|------|----------|-----------|
| `11111111-1111-1111-1111-111111111111` | 田中太郎 | 20230001 | `account` |
| `22222222-2222-2222-2222-222222222222` | 佐藤花子 | 20230002 | なし |
| `33333333-3333-3333-3333-333333333333` | 山田次郎 | 20230003 | `infra` |
| `44444444-4444-4444-4444-444444444444` | 鈴木美咲 | 20230004 | なし |
| `55555555-5555-5555-5555-555555555555` | 高橋健太 | 20230005 | なし |

---

## ローカル環境でのテスト手順

### 1. 開発環境の起動

```bash
make build && make up
make migrate
make insert_test
```

### 2. 認証の無効化

> **注意**: `AUTH=disable` は開発環境でのみ使用してください。本番環境では必ず有効な JWT 認証を使用してください。

`.env` に以下を設定します。

```
AUTH=disable
```

`AUTH=disable` 時の動作:
- JWT の署名検証はスキップ
- JWT の `sub` フィールドから `user_id` を取得する処理は有効
- OpenAPI スコープ（claim）の検証はバイパスされる

> claim 制御のテストは本番相当の設定（`AUTH` 有効）で行うか、コード内チェックのある `PUT /activity/record/{recordId}` で動作確認してください。

### 3. JWT トークンの生成

`AUTH=disable` でも `sub` を含む JWT 形式のトークンが必要です。`scripts/gen_jwt.sh` を使用して生成します。

```bash
# デフォルト（田中太郎: account ユーザー）
./scripts/gen_jwt.sh

# user_id を指定
./scripts/gen_jwt.sh <user_id>
```

**infra ユーザー（山田次郎）:**

```bash
./scripts/gen_jwt.sh 33333333-3333-3333-3333-333333333333
```

**account ユーザー（田中太郎）:**

```bash
./scripts/gen_jwt.sh 11111111-1111-1111-1111-111111111111
```

**一般ユーザー（佐藤花子）:**

```bash
./scripts/gen_jwt.sh 22222222-2222-2222-2222-222222222222
```

### 4. テストシナリオ

以下は `BASE_URL=http://localhost:8080` を前提とした curl の例です。

---

#### 管理者フラグの確認

```bash
# account ユーザー → isAdmin: true
curl -s -H "Authorization: Bearer <田中太郎のトークン>" \
  http://localhost:8080/user/me | jq .isAdmin

# 一般ユーザー → isAdmin: false
curl -s -H "Authorization: Bearer <佐藤花子のトークン>" \
  http://localhost:8080/user/me | jq .isAdmin
```

---

#### 通常グループの作成

```bash
curl -s -X POST http://localhost:8080/group \
  -H "Authorization: Bearer <佐藤花子のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"name":"テストグループ","description":"テスト用","joinable":true}'
```

---

#### 任意の claim を持つグループの作成（`POST /group/admin`）

```bash
# infra ユーザーで実行（AUTH=disable では一般ユーザーでも通過する）
curl -s -X POST http://localhost:8080/group/admin \
  -H "Authorization: Bearer <山田次郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"name":"新インフラチーム","description":"インフラ担当","joinable":false,"claim":"infra"}'
```

---

#### グループへのユーザー追加

```bash
GROUP_ID="<グループID>"

# グループ所属ユーザーとして追加
curl -s -X POST "http://localhost:8080/group/${GROUP_ID}/user" \
  -H "Authorization: Bearer <グループ作成者のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"userId":"33333333-3333-3333-3333-333333333333"}'

# infra ユーザーとして追加（グループ非所属でも可）
curl -s -X POST "http://localhost:8080/group/${GROUP_ID}/user" \
  -H "Authorization: Bearer <山田次郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"userId":"44444444-4444-4444-4444-444444444444"}'
```

---

#### グループへの自発参加

```bash
# joinable=true の通常グループに参加
curl -s -X POST "http://localhost:8080/group/f1111111-1111-1111-1111-111111111111/join" \
  -H "Authorization: Bearer <高橋健太のトークン>"

# 役職 claim（account/infra）を持つグループへの参加（403 が返る）
curl -s -X POST "http://localhost:8080/group/f0000007-1111-1111-1111-111111111111/join" \
  -H "Authorization: Bearer <佐藤花子のトークン>"
```

---

#### 管理者による強制チェックアウト

```bash
# 田中太郎をチェックイン
curl -s -X POST http://localhost:8080/activity/checkin \
  -H "Authorization: Bearer <田中太郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"place":"部室","checkInAt":"2026-02-18T10:00:00+09:00"}'

# infra ユーザーが強制チェックアウト
curl -s -X POST "http://localhost:8080/activity/checkout/11111111-1111-1111-1111-111111111111" \
  -H "Authorization: Bearer <山田次郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"place":"部室","checkoutAt":"2026-02-18T18:00:00+09:00"}'
```

---

#### 在室レコードの日時編集

```bash
RECORD_ID="<レコードID>"

# 自分のレコードを編集（所有者なので claim 不要）
curl -s -X PUT "http://localhost:8080/activity/record/${RECORD_ID}" \
  -H "Authorization: Bearer <田中太郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"activityType":"checkin","time":"2026-02-18T09:30:00+09:00"}'

# 他人のレコードを infra ユーザーが編集
curl -s -X PUT "http://localhost:8080/activity/record/${RECORD_ID}" \
  -H "Authorization: Bearer <山田次郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"activityType":"checkout","time":"2026-02-18T19:00:00+09:00"}'

# 一般ユーザーが他人のレコードを編集（403 が返る）
curl -s -X PUT "http://localhost:8080/activity/record/${RECORD_ID}" \
  -H "Authorization: Bearer <佐藤花子のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"activityType":"checkout","time":"2026-02-18T19:00:00+09:00"}'
```

---

#### 部費支払い一覧取得

```bash
# account ユーザーで実行
curl -s -X GET "http://localhost:8080/payment?year=2026" \
  -H "Authorization: Bearer <田中太郎のトークン>"
```

---

#### 稟議の承認

```bash
BUDGET_ID="a0000012-2222-2222-2222-222222222222"

curl -s -X PUT "http://localhost:8080/budget/${BUDGET_ID}/admin" \
  -H "Authorization: Bearer <田中太郎のトークン>" \
  -H "Content-Type: application/json" \
  -d '{"status":"approve","approverUserId":"11111111-1111-1111-1111-111111111111"}'
```

---

### 5. ログの確認

```bash
make logs
```

エラー時はサーバーログの `log` フィールドに詳細が出力されます。

### 6. データベースの直接確認

```bash
# グループと claim の確認
SELECT BIN_TO_UUID(g.id) AS group_id, g.name, gc.claim
FROM `groups` g
LEFT JOIN group_claims gc ON g.id = gc.group_id
ORDER BY g.name;

# ユーザーの役職 claim を確認
SELECT BIN_TO_UUID(gu.user_id) AS user_id, up.username, gc.claim
FROM groups_users gu
INNER JOIN group_claims gc ON gu.group_id = gc.group_id
INNER JOIN user_profiles up ON gu.user_id = up.user_id
ORDER BY gc.claim, up.username;
```
