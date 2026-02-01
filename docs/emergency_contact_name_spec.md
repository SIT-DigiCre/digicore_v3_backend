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
