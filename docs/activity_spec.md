# アクティビティ機能 仕様書

## 概要

部室への入退室管理機能を実装します。場所は`place`カラムで管理し、部室に限定しません。チェックインとチェックアウトは同じテーブル`activities`の1レコード内で`check_in_at`および`check_out_at`として管理し、最新レコードの`check_out_at`がNULLかどうかで在室状況を判定します。

## データベース設計

### テーブル: `activities`

```sql
CREATE TABLE activities
(
    id                     BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id                BINARY(16)   NOT NULL,
    place                  VARCHAR(255) NOT NULL,
    note                   TEXT         NULL,
    initial_checked_in_at  DATETIME     NOT NULL,
    initial_checked_out_at DATETIME     NULL,
    checked_in_at          DATETIME     NOT NULL,
    checked_out_at         DATETIME     NULL,
    created_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id_place_checked_in_at (user_id, place, checked_in_at DESC)
);
```

**カラム説明:**
- `id`: レコードID（UUID）
- `user_id`: ユーザーID
- `place`: 場所名（部室など）
- `note`: メモ（任意の補足情報、自由記述）
- `initial_checked_in_at`: 初回チェックイン日時（編集される前の値、常に不変）
- `initial_checked_out_at`: 初回チェックアウト日時（編集される前の値、常に不変、退室前はNULL）
- `checked_in_at`: 現在有効なチェックイン日時（編集により更新されうる）
- `checked_out_at`: 現在有効なチェックアウト日時（編集により更新されうる、退室前はNULL）
- `created_at`: レコード作成日時
- `updated_at`: レコード更新日時

**設計の利点:**
- 1レコードで入室〜退室までを一括管理できるため、チェックインとチェックアウトの対応関係が明確
- `initial_*`と`*At`を分離することで、編集履歴のトレースがしやすい
- 在室判定は最新のチェックインレコードの`check_out_at`がNULLかどうかで判定できるためシンプル

**在室判定ロジック:**
- ユーザーと場所でレコードを`checked_in_at`の降順でソート
- 最新レコードの`checked_out_at`がNULLの場合、在室中
- 最新レコードの`checked_out_at`がNULLでない場合、不在

## APIエンドポイント

### 1. チェックイン: `POST /activity/checkin`

**リクエスト:**
- `place`: 場所名（必須）
- `checkInAt`: チェックイン日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. ユーザーが既に同じ場所で在室中か確認（最新レコードの`checked_out_at`がNULLかどうか）
2. 在室中の場合：
   - 既存レコードの`initial_checked_out_at`および`checked_out_at`を、現在時刻または指定時刻で設定
   - 新しいチェックインレコードを作成
     - `initial_checked_in_at` = チェックイン日時
     - `checked_in_at` = チェックイン日時
     - `initial_checked_out_at` = NULL
     - `checked_out_at` = NULL
3. 在室中でない場合：
   - 新しいチェックインレコードを作成（上記と同様）

### 2. チェックアウト: `POST /activity/checkout`

**リクエスト:**
- `place`: 場所名（必須）
- `checkOutAt`: チェックアウト日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. ユーザーが指定場所で在室中か確認（最新レコードの`checked_out_at`がNULLかどうか）
2. 在室中の場合：
   - 対象レコードの`initial_checked_out_at`および`checked_out_at`を、現在時刻または指定時刻で設定
3. 在室中でない場合：
   - エラーを返す

### 3. 管理者による強制チェックアウト: `POST /activity/checkout/{userId}`

**リクエスト:**
- `place`: 場所名（必須）
- `checkOutAt`: チェックアウト日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. リクエスト送信者が管理者か確認（`pkg/admin/admin.go`の`CheckUserIsAdmin`を使用）
2. 対象ユーザーが指定場所で在室中か確認（最新レコードの`checked_out_at`がNULLかどうか）
3. 在室中の場合：
   - 対象レコードの`initial_checked_out_at`および`checked_out_at`を、現在時刻または指定時刻で設定

### 4. 現在在室中のユーザー一覧: `GET /activity/place/{place}/current`

**レスポンス:**
- `users`: ユーザー配列
  - `userId`: ユーザーID
  - `username`: ユーザー名
  - `shortIntroduction`: 自己紹介（短い）
  - `iconUrl`: アイコンURL
  - `checkedInAt`: 入室時刻

**処理フロー:**
1. 指定場所のレコードをユーザーごとに`checked_in_at`の降順でソート
2. 各ユーザーについて最新レコードの`checked_out_at`がNULLのものを抽出
3. 該当ユーザーのレコードの`checked_in_at`を入室時刻として取得
4. ユーザー情報（`user_profiles`）と結合
5. 入室時刻の昇順でソート

### 5. 過去の訪問履歴（ユーザーごとの入室回数）: `GET /activity/place/{place}/history`

**クエリパラメータ:**
- `period`: `day` | `week` | `month`（必須）
- `date`: 日付（YYYY-MM-DD形式、必須）

**レスポンス:**
- `users`: ユーザー配列
  - `userId`: ユーザーID
  - `username`: ユーザー名
  - `shortIntroduction`: 自己紹介（短い）
  - `iconUrl`: アイコンURL
  - `checkInCount`: 入室回数

**処理フロー:**
1. `period`に応じて日付範囲を計算
   - `day`: 指定日の00:00:00 ～ 23:59:59
   - `week`: 指定日を含む週の1週間前の月曜00:00:00 ～ 日曜23:59:59
   - `month`: 指定日を含む月の1か月前の1日00:00:00 ～ 月末23:59:59
2. 指定場所で`checked_in_at`が範囲内のレコードをユーザーごとに集計
3. 入室回数をカウント
4. ユーザー情報（`user_profiles`）と結合
5. 入室回数の降順でソート

### 6. ユーザーごとの入室記録: `GET /activity/user/{userId}/records`

**クエリパラメータ:**
- `place`: 場所名（オプション、指定しない場合は全場所）
- `offset`: オフセット（オプション、デフォルト0）
- `limit`: 取得件数（オプション、デフォルト50）

**レスポンス:**
- `records`: レコード配列
  - `recordId`: レコードID
  - `place`: 場所名
  - `checkedInAt`: 現在有効なチェックイン日時
  - `checkedOutAt`: 現在有効なチェックアウト日時（退室前はNULL）
  - `initialCheckedInAt`: 初回チェックイン日時
  - `initialCheckedOutAt`: 初回チェックアウト日時
- `total`: 総レコード数
- `offset`: 現在のオフセット
- `limit`: 現在のリミット

**処理フロー:**
1. 指定ユーザーのレコードを取得（`place`が指定されている場合はフィルタ）
2. `checked_in_at`の降順でソート
3. ページネーション適用（`offset`と`limit`）
4. 総レコード数を取得

### 7. レコード日時の更新: `PUT /activity/record/{recordId}`

**リクエスト:**
- `checkedInAt`: 編集後のチェックイン日時（オプション）
- `checkedOutAt`: 編集後のチェックアウト日時（オプション）

**処理フロー:**
1. レコードIDで`activities`レコードを取得
2. レコードが存在しない場合は404
3. リクエスト送信者がレコードの所有者か確認（管理者は除く）
4. `checkedInAt`が指定されている場合は`checked_in_at`のみを更新（`initial_checked_in_at`は変更しない）
5. `checkedOutAt`が指定されている場合は`checked_out_at`のみを更新（`initial_checked_out_at`は変更しない）
