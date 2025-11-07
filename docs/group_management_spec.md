# グループ管理機能 仕様書

## 概要

ユーザーがグループを作成し、他のユーザーをグループに追加できる機能を実装します。
また、管理者権限の判定機能を追加し、自分の情報を返すエンドポイントに管理者かどうかの情報を含めます。

## データベース構造

### 既存テーブル

#### `groups` テーブル

- グループの基本情報を格納
- 既存フィールド: id, name, description, joinable, user_count, created_at, updated_at

#### `groups_users` テーブル

- ユーザーとグループの関連を格納
- 既存フィールド: id, user_id, group_id, created_at, updated_at
- ユニーク制約: (user_id, group_id)

#### `group_claims` テーブル

- グループに対する権限(claim)を格納
- 既存フィールド: id, group_id, claim, created_at, updated_at
- ユニーク制約: (group_id, claim)
- claim 値の例: "admin"

## 機能要件

### 1. グループ作成機能

#### エンドポイント

```
POST /group
```

#### 権限

- 認証済みユーザー全員が通常のグループを作成可能
- `claim="admin"` のグループに所属しているユーザーのみが、`claim="admin"` のグループを作成可能

#### リクエスト

```json
{
  "name": "グループ名",
  "description": "グループの説明",
  "joinable": true,
  "isAdminGroup": false
}
```

#### レスポンス

```json
{
  "groupId": "uuid",
  "name": "グループ名",
  "description": "グループの説明",
  "joinable": true,
  "userCount": 1
}
```

#### 処理フロー

1. ユーザー認証を確認
2. `isAdminGroup=true` の場合、リクエストユーザーが管理者かどうかを確認
   - 管理者でない場合は 403 Forbidden を返す
3. グループを作成 (`groups` テーブルに INSERT)
4. 作成者をグループメンバーとして追加 (`groups_users` テーブルに INSERT)
5. `isAdminGroup=true` の場合、`group_claims` テーブルに `claim="admin"` を追加
6. 作成したグループ情報を返す

### 2. グループメンバー追加機能

#### エンドポイント

```
POST /group/{group_id}/user
```

#### 権限

- グループに既に所属しているユーザーのみが実行可能

#### リクエスト

```json
{
  "userId": "追加するユーザーのuuid"
}
```

#### レスポンス

```json
{
  "message": "ユーザーをグループに追加しました"
}
```

#### 処理フロー

1. ユーザー認証を確認
2. リクエストユーザーがグループに所属しているか確認
   - 所属していない場合は 403 Forbidden を返す
3. 追加対象ユーザーが存在するか確認
   - 存在しない場合は 404 Not Found を返す
4. 追加対象ユーザーが既にグループに所属していないか確認
   - 既に所属している場合は 400 Bad Request を返す
5. `groups_users` テーブルに追加
6. `groups` テーブルの `user_count` をインクリメント
7. 成功メッセージを返す

### 3. グループ自発参加機能

#### エンドポイント

```
POST /group/{group_id}/join
```

#### 権限

- 認証済みユーザーが実行可能
- `joinable=true` のグループのみ参加可能
- `claim="admin"` のグループには参加不可（`joinable` の値に関わらず）

#### リクエスト

なし（リクエストボディ不要）

#### レスポンス

```json
{
  "message": "グループに参加しました"
}
```

#### 処理フロー

1. ユーザー認証を確認
2. グループが存在するか確認
   - 存在しない場合は 404 Not Found を返す
3. グループが `claim="admin"` かどうか確認
   - admin グループの場合は 403 Forbidden を返す（エラーメッセージ: "このグループには参加できません"）
4. グループの `joinable` が `true` かどうか確認
   - `false` の場合は 403 Forbidden を返す（エラーメッセージ: "このグループには参加できません"）
5. リクエストユーザーが既にグループに所属していないか確認
   - 既に所属している場合は 400 Bad Request を返す（エラーメッセージ: "既にグループに参加しています"）
6. `groups_users` テーブルに追加
7. `groups` テーブルの `user_count` をインクリメント
8. 成功メッセージを返す

### 4. 自分の情報取得機能の拡張

#### エンドポイント

```
GET /user/me
```

#### 既存レスポンスへの追加

```json
{
  "userId": "uuid",
  "username": "ユーザー名",
  "studentNumber": "学籍番号",
  "iconUrl": "アイコンURL",
  "schoolGrade": 1,
  "discordUserId": "Discord ID",
  "shortIntroduction": "自己紹介",
  "activeLimit": "有効期限",
  "isAdmin": true
}
```

#### 管理者判定ロジック

- `group_claims` テーブルで `claim="admin"` のグループを取得
- そのグループに `groups_users` テーブルでユーザーが所属しているかを確認
- 1 つでも所属していれば `isAdmin=true`

#### SQL 例

```sql
SELECT COUNT(*) > 0 AS is_admin
FROM groups_users
INNER JOIN group_claims ON groups_users.group_id = group_claims.group_id
WHERE groups_users.user_id = UUID_TO_BIN(?)
  AND group_claims.claim = 'admin'
```

## 設計上の決定事項

### 権限管理

- **グループメンバー追加**: グループに所属している任意のユーザーが他のユーザーを追加可能
- **グループ自発参加**: `joinable=true` の一般グループのみ可能。管理者グループ（`claim="admin"`）には不可

### 実装しない機能

以下の機能は今回の実装では含めません:

- グループからのメンバー削除機能
- グループ削除機能
- グループ編集機能（名前、説明の変更）
- グループ参加・追加時の通知機能

### エラーメッセージ

- 全てのエラーメッセージは日本語で記述
- 既存のパターン（`pkg/api/response`）に従う

### 拡張性

- `group_claims` テーブルは将来的に "admin" 以外の claim タイプ（例: "moderator", "member" など）を追加できる設計とする

## 実装ファイル予定

### API 定義 (OpenAPI)

- `document/paths/group.yml` - POST /group の追加
- `document/paths/group_group_id_user.yml` - 新規作成（POST /group/{group_id}/user）
- `document/paths/group_group_id_join.yml` - 新規作成（POST /group/{group_id}/join）
- `document/paths/user_me.yml` - レスポンス更新は不要（スキーマのみ更新）
- `document/schemas/req_post_group.yml` - 新規作成
- `document/schemas/res_post_group.yml` - 新規作成
- `document/schemas/req_post_group_group_id_user.yml` - 新規作成
- `document/schemas/res_get_user_me.yml` - isAdmin フィールド追加

### Go 実装

- `pkg/group/post_group.go` - 新規作成
- `pkg/group/post_group_group_id_user.go` - 新規作成
- `pkg/group/post_group_group_id_join.go` - 新規作成
- `pkg/user/get_user_me.go` - 管理者判定ロジック追加

### SQL

- `pkg/db/sql/group/insert_group.sql` - 新規作成
- `pkg/db/sql/group/insert_groups_users.sql` - 新規作成
- `pkg/db/sql/group/insert_group_claims.sql` - 新規作成
- `pkg/db/sql/group/select_user_is_admin.sql` - 新規作成
- `pkg/db/sql/group/select_is_group_member.sql` - 新規作成
- `pkg/db/sql/group/select_group_is_admin_group.sql` - 新規作成（グループが admin グループか確認）
- `pkg/db/sql/group/select_group_joinable.sql` - 新規作成（グループの joinable 状態を取得）
- `pkg/db/sql/group/update_group_user_count_increment.sql` - 新規作成

## 備考

- トランザクション管理が必要な処理（グループ作成、メンバー追加、グループ参加）では `TransactionClient` を使用します
- エラーハンドリングは既存の `pkg/api/response` パッケージのパターンに従います
- 全てのエラーメッセージは日本語で記述します
- グループ参加機能では、管理者グループへの参加を防ぐため、`group_claims` テーブルの確認を必ず行います

## Postman での動作確認手順

### 環境準備

#### 開発環境の起動

```bash
# プロジェクトルートで実行
make build && make up
make migrate
make insert_test
```

#### 認証の無効化設定

`.env`ファイルで以下を設定：

```
AUTH=disable
```

### Postman コレクション設定

#### 認証ヘッダー設定

認証が無効化されていても、JWT トークンの subject から user_id を取得するため、有効な JWT 形式のトークンが必要です。以下のようなダミートークンを生成：

```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTEiLCJleHAiOjk5OTk5OTk5OTl9.dummy
```

トークンの生成は https://www.jwt.io/ja から行ってください。

#### ヘッダーの例

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

#### ペイロードの例

```json
{
  "sub": "55555555-5555-5555-5555-555555555555"
}
```

### テストシナリオ

#### 1. 事前確認：自分の情報取得

```
GET {{base_url}}/user/me
Authorization: Bearer [上記のトークン]
```

**期待結果**: `isAdmin: false` が含まれる

#### 2. 通常グループの作成

```
POST {{base_url}}/group
Content-Type: application/json
Authorization: Bearer [user_id_1のトークン]

{
  "name": "テストグループ1",
  "description": "Postmanテスト用のグループです",
  "joinable": true,
  "isAdminGroup": false
}
```

**期待結果**:

- ステータス: 200
- レスポンスに `groupId`, `userCount: 1` が含まれる

#### 3. 管理者グループの作成（権限不足でエラー）

```
POST {{base_url}}/group
Content-Type: application/json
Authorization: Bearer [user_id_1のトークン]

{
  "name": "管理者グループ",
  "description": "管理者専用グループ",
  "joinable": false,
  "isAdminGroup": true
}
```

**期待結果**:

- ステータス: 403
- エラーメッセージ: "管理者グループを作成する権限がありません"

#### 4. グループメンバー追加

```
POST {{base_url}}/group/{{group_id}}/user
Content-Type: application/json
Authorization: Bearer [user_id_1のトークン]

{
  "userId": "22222222-2222-2222-2222-222222222222"
}
```

**期待結果**:

- ステータス: 200
- メッセージ: "ユーザーをグループに追加しました"

#### 5. グループ自発参加

```
POST {{base_url}}/group/{{group_id}}/join
Authorization: Bearer [user_id_3のトークン]
```

**期待結果**:

- ステータス: 200
- メッセージ: "グループに参加しました"

#### 6. 既存グループへの重複参加（エラー）

```
POST {{base_url}}/group/{{group_id}}/join
Authorization: Bearer [user_id_3のトークン]
```

**期待結果**:

- ステータス: 400
- エラーメッセージ: "既にグループに参加しています"

### 管理者権限のテスト

#### 管理者グループの手動作成

データベースに直接管理者グループを作成：

```sql
-- 管理者グループ作成
INSERT INTO `groups` (id, name, description, joinable, user_count) VALUES
  (UUID_TO_BIN('admin111-1111-1111-1111-111111111111'), '管理者グループ', '管理者専用', false, 1);

-- 管理者権限付与
INSERT INTO group_claims (id, group_id, claim) VALUES
  (UUID_TO_BIN('claim111-1111-1111-1111-111111111111'), UUID_TO_BIN('admin111-1111-1111-1111-111111111111'), 'admin');

-- ユーザーを管理者グループに追加
INSERT INTO groups_users (id, group_id, user_id) VALUES
  (UUID_TO_BIN('gu111111-1111-1111-1111-111111111111'), UUID_TO_BIN('admin111-1111-1111-1111-111111111111'), UUID_TO_BIN('11111111-1111-1111-1111-111111111111'));
```

#### 管理者権限確認

```
GET {{base_url}}/user/me
Authorization: Bearer [user_id_1のトークン]
```

**期待結果**: `isAdmin: true` が含まれる

#### 管理者による管理者グループ作成

```
POST {{base_url}}/group
Content-Type: application/json
Authorization: Bearer [user_id_1のトークン]

{
  "name": "新しい管理者グループ",
  "description": "管理者が作成した管理者グループ",
  "joinable": false,
  "isAdminGroup": true
}
```

**期待結果**:

- ステータス: 200
- 正常にグループが作成される

### エラーケースのテスト

#### 存在しないユーザーの追加

```
POST {{base_url}}/group/{{group_id}}/user
Content-Type: application/json
Authorization: Bearer [user_id_1のトークン]

{
  "userId": "99999999-9999-9999-9999-999999999999"
}
```

**期待結果**:

- ステータス: 404
- エラーメッセージ: "指定されたユーザーが存在しません"

#### 存在しないグループへの参加

```
POST {{base_url}}/group/99999999-9999-9999-9999-999999999999/join
Authorization: Bearer [user_id_1のトークン]
```

**期待結果**:

- ステータス: 404
- エラーメッセージ: "指定されたグループが存在しません"

#### 管理者グループへの自発参加（エラー）

```
POST {{base_url}}/group/{{admin_group_id}}/join
Authorization: Bearer [user_id_2のトークン]
```

**期待結果**:

- ステータス: 403
- エラーメッセージ: "このグループには参加できません"

### テスト時の注意事項

1. **認証トークン**: `AUTH=disable`設定時でも、JWT トークンの subject から user_id を取得するため、有効な JWT 形式のトークンが必要です
2. **データベース確認**: 各操作後にデータベースを直接確認して、データが正しく挿入されているか確認してください
3. **ログ確認**: `make logs`でサーバーログを確認し、エラーの詳細を把握してください
4. **グループ ID**: テスト中に作成されたグループの ID は、レスポンスから取得して後続のテストで使用してください
