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
- `pkg/db/sql/group/select_group_is_admin_group.sql` - 新規作成（グループがadminグループか確認）
- `pkg/db/sql/group/select_group_joinable.sql` - 新規作成（グループのjoinable状態を取得）
- `pkg/db/sql/group/update_group_user_count_increment.sql` - 新規作成

## 備考

- トランザクション管理が必要な処理（グループ作成、メンバー追加、グループ参加）では `TransactionClient` を使用します
- エラーハンドリングは既存の `pkg/api/response` パッケージのパターンに従います
- 全てのエラーメッセージは日本語で記述します
- グループ参加機能では、管理者グループへの参加を防ぐため、`group_claims` テーブルの確認を必ず行います
