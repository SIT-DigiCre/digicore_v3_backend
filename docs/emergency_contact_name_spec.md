# 緊急連絡先氏名フィールド分割 仕様書

## 概要

緊急連絡先の氏名を「名字」「名前」の2フィールドで扱う。データベースでは `parent_last_name`（名字）および `parent_first_name`（名前）カラムを使用する。

旧フィールド `parent_name`（氏名フルネーム）は DB カラムとしては残存するが（`DEFAULT ''`）、API・SQL からは参照・更新しない。

- **名字**: `parent_last_name` カラム。姓のみを格納する。
- **名前**: `parent_first_name` カラム。名のみを格納する。

本人の氏名と同様、緊急連絡先も **parentLastName = 名字（姓）、parentFirstName = 名前（名）** の命名とする。

## データベース

### 対象テーブル: `user_private_profiles`

#### 使用カラム

- `parent_last_name VARCHAR(255) NOT NULL DEFAULT ''` … 緊急連絡先の名字
- `parent_first_name VARCHAR(255) NOT NULL DEFAULT ''` … 緊急連絡先の名前

#### 廃止済みカラム（残存）

- `parent_name VARCHAR(255) NOT NULL DEFAULT ''` … 旧・緊急連絡先氏名。API からは参照・更新しない。既存データの保持のみ。

## API レスポンス（取得）

### `GET /user/me/private`

レスポンスに以下を含める。

| フィールド | 説明 | 対応DBカラム |
|-----------|------|--------------|
| `parentLastName` | 緊急連絡先の名字 | `parent_last_name` |
| `parentFirstName` | 緊急連絡先の名前 | `parent_first_name` |

## API リクエスト（作成・更新）

### `PUT /user/me/private`

リクエストボディの緊急連絡先氏名まわりは次のようにマッピングする。

| フィールド | 説明 | 更新先DBカラム |
|-----------|------|----------------|
| `parentLastName` | 緊急連絡先の名字 | `parent_last_name` |
| `parentFirstName` | 緊急連絡先の名前 | `parent_first_name` |

- OpenAPI 上では `parentLastName` / `parentFirstName` を **required にしない**。
- 未送信のフィールドは更新しない（既存値維持）。部分更新として扱う。
- `parentLastName` と `parentFirstName` は必ずセットで送信する必要がある。片方のみの送信は 400 Bad Request を返す。

## バリデーション

- **送信したフィールドに空文字を許容しない**: リクエストに含まれた `parentLastName` / `parentFirstName` のいずれかが空文字の場合は、400 Bad Request を返す。
- 未送信のフィールドは「空文字」としては扱わず、更新対象から外す（既存値維持）。

## 後方互換性

- 旧 `parentName` フィールドは API から完全に削除済み。リクエストに含まれていても Go の JSON デコーダにより無視される。
- DB の `parent_name` カラムは `DEFAULT ''` で残存し、既存データは保持される。

## 実装時の変更ファイル一覧（参考）

| 種別 | ファイル |
|------|----------|
| OpenAPI | `document/schemas/req_put_user_me_private.yml` … `parentLastName`, `parentFirstName` のみ |
| OpenAPI | `document/schemas/res_get_user_me_private.yml` … `parentLastName`, `parentFirstName` のみ |
| SQL | `pkg/db/sql/user/select_user_private_from_user_id.sql` … `parent_last_name`, `parent_first_name` を SELECT |
| SQL | `pkg/db/sql/user/update_user_private.sql` … `parent_last_name`, `parent_first_name` を UPDATE |
| SQL | `pkg/db/sql/user/insert_user_private.sql` … `parent_last_name`, `parent_first_name` を INSERT |
| SQL | `pkg/db/sql/user/insert_user_private_default.sql` … デフォルト INSERT（`parent_name` は含めない） |
| Go | `pkg/user/get_user_me_private.go` … private 構造体に `ParentLastName`, `ParentFirstName` のみ |
| Go | `pkg/user/put_user_me_private.go` … `parentLastName` / `parentFirstName` のペアバリデーションと更新処理 |
| 生成 | `make generate_api` 実行により `pkg/api/models.gen.go` 等の生成コードを更新 |

※ `*.gen.go` および `*.gen.yml` は直接編集せず、OpenAPI 修正後に `make generate_api` で更新する。

---

## テスト項目

以下は `GET /user/me/private` および `PUT /user/me/private` の動作確認用テスト項目である。認証は Bearer トークン（JWT の `sub` にユーザーID）が必要。開発環境では `.env` で `AUTH=disable` にし、署名検証なしでトークンを扱える。

### 前提

- ベースURL: `http://localhost:8000`（`BACKEND_PORT` に合わせる）
- 認証: `Authorization: Bearer <JWT>`（`sub` に対象ユーザーの UUID）
- 対象ユーザーに `user_private_profiles` が存在すること（存在しない場合は PUT で新規作成される）
- **テスト実行前**: (1) `make migrate` を実行し、`user_private_profiles` に `parent_last_name` / `parent_first_name` が存在することを確認する。(2) 実装反映後はバックエンドを再ビルド・再起動（例: `make build && make up`）してから叩くこと。

### テスト一覧

| No. | エンドポイント | 内容 | 期待結果 |
|-----|----------------|------|----------|
| 1 | GET /user/me/private | 個人情報取得 | 200。レスポンスに `parentLastName`, `parentFirstName` が含まれる。 |
| 2 | PUT /user/me/private | parentLastName / parentFirstName を指定して更新 | 200。レスポンスで同じ値が返る。 |
| 3 | GET /user/me/private | 更新後の取得 | 200。No.2 で送った `parentLastName`, `parentFirstName` がそのまま返る。 |
| 4 | PUT /user/me/private | parentLastName / parentFirstName を送らずに更新 | 200。名字・名前は既存値が維持される。 |
| 5 | PUT /user/me/private | parentLastName のみ送信（parentFirstName なし） | 400 Bad Request。名字と名前は両方指定が必要。 |
| 6 | PUT /user/me/private | parentLastName に空文字を送信 | 400 Bad Request。「緊急連絡先の名字に空文字は指定できません」の旨。 |
| 7 | PUT /user/me/private | parentFirstName に空文字を送信 | 400 Bad Request。「緊急連絡先の名前に空文字は指定できません」の旨。 |

### curl 例（テスト用）

```bash
# テスト用 JWT（make jwt で生成するか、以下をコピー）
TOKEN="$(make jwt | head -1)"

# No.1 GET 取得
curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" http://localhost:8000/user/me/private

# No.2 PUT 名字・名前指定
curl -s -w "\n%{http_code}" -X PUT -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"lastName":"山田","lastNameKana":"ヤマダ","firstName":"太郎","firstNameKana":"タロウ","isMale":true,"phoneNumber":"09012345678","address":"東京都","parentLastName":"山田","parentFirstName":"花子","parentCellphoneNumber":"09098765432","parentHomephoneNumber":"0312345678","parentAddress":"東京都"}' \
  http://localhost:8000/user/me/private
```

（PUT のボディは必須項目をすべて含む。`parentLastName` / `parentFirstName` は省略可能だが、送る場合は両方セットで。`parentHomephoneNumber` を送る場合は電話番号形式であること。）
