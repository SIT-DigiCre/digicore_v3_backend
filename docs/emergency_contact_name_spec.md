# 緊急連絡先氏名フィールド分割 仕様書

## 概要

緊急連絡先の氏名を「氏名」「名字」「名前」の3フィールドで扱い、データベースでは既存の氏名カラムに加えて名字・名前用のカラムを追加します。既存の `parent_name`（氏名）カラムは残し、`parent_last_name`（名字）および `parent_first_name`（名前）を新規追加することで、既存クライアントの動作を維持しつつ、名字・名前を個別に扱えるようにします。

- **氏名**: 既存の `parent_name` カラム。フルネームとしての表示・更新に利用する。
- **名字**: 新規 `parent_last_name` カラム。姓のみを格納する。
- **名前**: 新規 `parent_first_name` カラム。名のみを格納する。

本人の氏名と同様、緊急連絡先も **parentLastName = 名字（姓）、parentFirstName = 名前（名）** の命名とします。

## データベース

### 対象テーブル: `user_private_profiles`

#### 既存カラム（変更なし）

- `parent_name VARCHAR(255) NOT NULL` … 緊急連絡先氏名（フルネーム）

#### 追加カラム

- `parent_last_name VARCHAR(255) NOT NULL DEFAULT ''` … 緊急連絡先の名字
- `parent_first_name VARCHAR(255) NOT NULL DEFAULT ''` … 緊急連絡先の名前

#### マイグレーション

- 上記2カラムを `ALTER TABLE` で追加するのみとする。
- 既存行へのデータ投入は行わない。既存レコードは `parent_last_name` および `parent_first_name` がデフォルトの空文字のままとなる。

## API レスポンス（取得）

### `GET /user/me/private`

レスポンスに以下を含める。

| フィールド | 説明 | 対応DBカラム |
|-----------|------|--------------|
| `parentName` | 緊急連絡先氏名（既存） | `parent_name` |
| `parentLastName` | 緊急連絡先の名字（新規） | `parent_last_name` |
| `parentFirstName` | 緊急連絡先の名前（新規） | `parent_first_name` |

既存クライアントは `parentName` のみを参照しても従来どおり動作する。新規クライアントは `parentLastName` / `parentFirstName` も利用可能。

## API リクエスト（作成・更新）

### `PUT /user/me/private`

リクエストボディの緊急連絡先氏名まわりは次のようにマッピングする。

| フィールド | 説明 | 更新先DBカラム |
|-----------|------|----------------|
| `parentName` | 緊急連絡先氏名 | `parent_name` |
| `parentLastName` | 緊急連絡先の名字 | `parent_last_name` |
| `parentFirstName` | 緊急連絡先の名前 | `parent_first_name` |

- OpenAPI 上では `parentName` / `parentLastName` / `parentFirstName` を **required にしない**。既存クライアントが `parentName` のみを送る利用形態を許容する。
- 未送信のフィールドは更新しない（既存値維持）。部分更新として扱う。

## バリデーション

- **送信したフィールドに空文字を許容しない**: リクエストに含まれた `parentName` / `parentLastName` / `parentFirstName` のいずれかが空文字の場合は、400 Bad Request を返す。
- 未送信のフィールドは「空文字」としては扱わず、更新対象から外す（既存値維持）。このため、従来どおり `parentName` のみを送るクライアントはそのまま動作する。

## 後方互換性

- **既存クライアント**: `parentName` のみの送受信で従来どおり利用可能。レスポンスに `parentLastName` / `parentFirstName` が追加されても、無視すればよい。
- **新規クライアント**: `parentLastName` / `parentFirstName` を送受信することで、名字・名前を個別に扱える。

## 実装時の変更ファイル一覧（参考）

| 種別 | ファイル |
|------|----------|
| スキーマ | 新規マイグレーション（`parent_last_name`, `parent_first_name` の追加） |
| OpenAPI | `document/schemas/req_put_user_me_private.yml` … プロパティ `parentLastName`, `parentFirstName` を追加（required には含めない） |
| OpenAPI | `document/schemas/res_get_user_me_private.yml` … プロパティ `parentLastName`, `parentFirstName` を追加 |
| SQL | `pkg/db/sql/user/select_user_private_from_user_id.sql` … 新カラムを SELECT に含める |
| SQL | `pkg/db/sql/user/update_user_private.sql` … 新カラムを UPDATE に含める |
| SQL | `pkg/db/sql/user/insert_user_private.sql` … 新カラムを INSERT に含める |
| Go | `pkg/user/get_user_me_private.go` … private 構造体およびレスポンスマッピングに新フィールドを追加 |
| Go | `pkg/user/put_user_me_private.go` … パラメータ構造体に新フィールドを追加し、空文字チェックのバリデーションを実装 |
| 生成 | `make generate_api` 実行により `pkg/api/models.gen.go` 等の生成コードを更新 |

※ `*.gen.go` および `*.gen.yml` は直接編集せず、OpenAPI 修正後に `make generate_api` で更新する。

---

## テスト項目

以下は `GET /user/me/private` および `PUT /user/me/private` の動作確認用テスト項目である。認証は Bearer トークン（JWT の `sub` にユーザーID）が必要。開発環境では `.env` で `AUTH=disable` にし、署名検証なしでトークンを扱える。

**JWT（テスト用トークン）**: プロジェクトルートで `make jwt` を実行すると有効なテスト用 JWT が出力される（`scripts/gen_jwt.sh`）。第1引数で user_id を指定可能。形式の詳細は [グループ管理機能 仕様書](docs/group_management_spec.md) の「認証ヘッダー設定」を参照。

### 前提

- ベースURL: `http://localhost:8000`（`BACKEND_PORT` に合わせる）
- 認証: `Authorization: Bearer <JWT>`（`sub` に対象ユーザーの UUID）。形式は上記グループ管理仕様書に同じ。
- 対象ユーザーに `user_private_profiles` が存在すること（存在しない場合は PUT で新規作成される）
- **テスト実行前**: (1) `make migrate` を実行し、`user_private_profiles` に `parent_last_name` / `parent_first_name` が存在することを確認する。(2) 実装反映後はバックエンドを再ビルド・再起動（例: `make build && make up`）してから叩くこと。

### テスト一覧

| No. | エンドポイント | 内容 | 期待結果 |
|-----|----------------|------|----------|
| 1 | GET /user/me/private | 個人情報取得 | 200。レスポンスに `parentName`, `parentLastName`, `parentFirstName` が含まれる。既存データのみの場合は `parentLastName` / `parentFirstName` は空文字の可能性あり。 |
| 2 | PUT /user/me/private | 3フィールドすべて指定して更新 | 200。`parentName`, `parentLastName`, `parentFirstName` をすべて送信。レスポンスで同じ値が返る。 |
| 3 | GET /user/me/private | 更新後の取得 | 200。No.2 で送った `parentName`, `parentLastName`, `parentFirstName` がそのまま返る。 |
| 4 | PUT /user/me/private | 既存クライアント様式（parentName のみ送信） | 200。`parentName` のみ含め、`parentLastName` / `parentFirstName` は送らない。既存の名字・名前は維持され、氏名のみ更新される。 |
| 5 | GET /user/me/private | 部分更新後の取得 | 200。No.4 で更新した `parentName` が反映され、`parentLastName` / `parentFirstName` は前回（No.2）の値のまま。 |
| 6 | PUT /user/me/private | parentLastName / parentFirstName のみ送信 | 200。名字・名前のみ送信し、`parentName` は送らない。氏名は既存値維持、名字・名前のみ更新。 |
| 7 | PUT /user/me/private | parentName に空文字を送信 | 400 Bad Request。メッセージに「緊急連絡先氏名に空文字は指定できません」等の旨が含まれる。 |
| 8 | PUT /user/me/private | parentLastName に空文字を送信 | 400 Bad Request。メッセージに「緊急連絡先の名字に空文字は指定できません」等の旨が含まれる。 |
| 9 | PUT /user/me/private | parentFirstName に空文字を送信 | 400 Bad Request。メッセージに「緊急連絡先の名前に空文字は指定できません」等の旨が含まれる。 |

### curl 例（テスト用）

```bash
# テスト用 JWT（make jwt で生成するか、以下をコピー）
TOKEN="$(make jwt | head -1)"

# No.1 GET 取得
curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" http://localhost:8000/user/me/private

# No.2 PUT 3フィールド指定（parentHomephoneNumber は省略可。送る場合は電話番号形式）
curl -s -w "\n%{http_code}" -X PUT -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"lastName":"山田","lastNameKana":"ヤマダ","firstName":"太郎","firstNameKana":"タロウ","isMale":true,"phoneNumber":"09012345678","address":"東京都","parentName":"山田花子","parentLastName":"山田","parentFirstName":"花子","parentCellphoneNumber":"09098765432","parentHomephoneNumber":"0312345678","parentAddress":"東京都"}' \
  http://localhost:8000/user/me/private
```

（PUT のボディは必須項目をすべて含む。緊急連絡先氏名まわりは `parentName` / `parentLastName` / `parentFirstName` を必要に応じて省略可能。`parentHomephoneNumber` を送る場合は電話番号形式であること。）

### 実施結果の記録（任意）

| No. | 期待HTTP | 実際のHTTP | 備考 |
|-----|----------|------------|------|
| 1 | 200 | 200 または 404 | 404 は個人情報未登録時。PUT で作成後に 200 |
| 2 | 200 | **200** | parentName/parentLastName/parentFirstName 反映 |
| 3 | 200 | **200** | No.2 の値が返る |
| 4 | 200 | **200** | parentName のみ更新、名字・名前は維持 |
| 5 | 200 | **200** | 部分更新の確認 |
| 6 | 200 | **200** | parentLastName/parentFirstName のみ更新 |
| 7 | 400 | **400** | 緊急連絡先氏名の長さは少なくとも1文字は… |
| 8 | 400 | **400** | 緊急連絡先の名字の長さは少なくとも1文字は… |
| 9 | 400 | **400** | 緊急連絡先の名前の長さは少なくとも1文字は… |

- **JWT**: テスト用トークンの形式は [group_management_spec.md](docs/group_management_spec.md) の「認証ヘッダー設定」に従うこと（jwt.io で `sub` にユーザー UUID、必要なら `exp` を含めたダミートークンを生成する）。
- **「不明なエラーが発生しました」の原因**: このメッセージは `get_user_me_private.go` では `dbClient.Select()` が `err != nil` のとき、`put_user_me_private.go` では `dbClient.DuplicateUpdate()` が `err != nil` のときに返る。**実際の原因はクライアントには返さず、サーバー側のログ（`logrus.Error(err.Log)`）にのみ出力される**。500 が出たらバックエンドのコンテナログ（`make logs` や `docker compose logs backend`）で `err.Log` の内容を確認すること。想定される原因例: (1) テスト用 JWT が不正で `user_id` が期待どおりでない、(2) `make migrate` 前で DB に `parent_last_name` / `parent_first_name` が無い、(3) SQL や struct の不整合。
