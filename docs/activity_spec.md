# アクティビティ機能 仕様書

## 概要

部室への入退室管理機能を実装します。場所は`place`カラムで管理し、部室に限定しません。チェックインとチェックアウトを同じテーブルで`type`カラム（VARCHAR）で判別し、最新レコードの`type`で在室状況を判定します。

## データベース設計

### テーブル: `activity_records`

```sql
CREATE TABLE activity_records
(
    id              BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id         BINARY(16)   NOT NULL,
    place           VARCHAR(255) NOT NULL,
    type            VARCHAR(255) NOT NULL,
    datetime        DATETIME     NOT NULL,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id_place_datetime (user_id, place, datetime DESC)
);
```

**カラム説明:**
- `id`: レコードID（UUID）
- `user_id`: ユーザーID
- `place`: 場所名（部室など）
- `type`: レコードタイプ（`checkin`、`checkout`など、VARCHARで拡張可能）
- `datetime`: チェックイン/チェックアウト日時（後から変更可能）
- `created_at`: レコード作成日時（自動）
- `updated_at`: レコード更新日時（自動）

**設計の利点:**
- `type`をVARCHARにすることで、将来的に他のタイプ（例：`break`、`return`など）を追加しやすい
- `related_id`がないため、シンプルな構造
- 在室判定は最新レコードの`type`を確認するだけなのでシンプル

**在室判定ロジック:**
- ユーザーと場所でレコードを`datetime`の降順でソート
- 最新レコードの`type`が`checkin`の場合、在室中
- 最新レコードの`type`が`checkout`の場合、不在

## APIエンドポイント

### 1. チェックイン: `POST /activity/checkin`

**リクエスト:**
- `place`: 場所名（必須）
- `checkInAt`: チェックイン日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. ユーザーが既に同じ場所でチェックイン中か確認（最新レコードの`type`が`checkin`かどうか）
2. チェックイン中の場合：
   - チェックアウトレコードを作成（`type='checkout'`、`datetime`は現在時刻または指定時刻）
   - 新しいチェックインレコードを作成（`type='checkin'`）
3. チェックイン中でない場合：
   - 新しいチェックインレコードを作成（`type='checkin'`）

### 2. チェックアウト: `POST /activity/checkout`

**リクエスト:**
- `place`: 場所名（必須）
- `checkOutAt`: チェックアウト日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. ユーザーが指定場所でチェックイン中か確認（最新レコードの`type`が`checkin`かどうか）
2. チェックイン中の場合：
   - チェックアウトレコードを作成（`type='checkout'`）
3. チェックイン中でない場合：
   - エラーを返す

### 3. 管理者による強制チェックアウト: `POST /activity/checkout/{userId}`

**リクエスト:**
- `place`: 場所名（必須）
- `checkOutAt`: チェックアウト日時（オプション、未指定時は現在時刻）

**処理フロー:**
1. リクエスト送信者が管理者か確認（`pkg/group/group.go`の`checkUserIsAdmin`を使用）
2. 対象ユーザーが指定場所でチェックイン中か確認（最新レコードの`type`が`checkin`かどうか）
3. チェックイン中の場合：
   - チェックアウトレコードを作成（`type='checkout'`）

### 4. 現在在室中のユーザー一覧: `GET /activity/place/{place}/current`

**レスポンス:**
- `users`: ユーザー配列
  - `userId`: ユーザーID
  - `username`: ユーザー名
  - `shortIntroduction`: 自己紹介（短い）
  - `iconUrl`: アイコンURL
  - `checkInAt`: 入室時刻

**処理フロー:**
1. 指定場所のレコードをユーザーごとに`datetime`の降順でソート
2. 各ユーザーについて最新レコードの`type`が`checkin`のものを抽出
3. 該当ユーザーの最新の`type='checkin'`レコードの`datetime`を取得
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
   - `week`: 指定日を含む週の月曜00:00:00 ～ 日曜23:59:59
   - `month`: 指定日を含む月の1日00:00:00 ～ 月末23:59:59
2. 指定場所で`type='checkin'`かつ`datetime`が範囲内のレコードをユーザーごとに集計
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
  - `type`: レコードタイプ（`checkin`または`checkout`）
  - `datetime`: 日時
  - `checkOutAt`: 退室時刻（`type='checkin'`の場合のみ、次の`checkout`レコードの`datetime`）
- `total`: 総レコード数
- `offset`: 現在のオフセット
- `limit`: 現在のリミット

**処理フロー:**
1. 指定ユーザーのレコードを取得（`place`が指定されている場合はフィルタ）
2. `datetime`の降順でソート
3. ページネーション適用（`offset`と`limit`）
4. チェックインレコードの場合、その後の最初のチェックアウトレコードを探して`checkOutAt`を設定
5. 総レコード数を取得

### 7. レコード日時の更新: `PUT /activity/record/{recordId}`

**リクエスト:**
- `datetime`: チェックイン/チェックアウト日時（必須）

**処理フロー:**
1. レコードIDで`activity_records`レコードを取得
2. レコードが存在しない場合は404
3. リクエスト送信者がレコードの所有者か確認（管理者は除く）
4. 指定された日時を更新
5. **注意**: `datetime`を変更すると在室判定に影響する可能性があるため、変更後の整合性チェックが必要な場合がある

## 注意事項

- すべてのエンドポイントは`BearerAuth`セキュリティスキームを要求
- 管理者権限チェックは`pkg/group/group.go`の`checkUserIsAdmin`関数を使用
- 日時は`DATETIME`型で保存し、タイムゾーンは`Asia/Tokyo`を使用
- 在室判定は、ユーザーと場所でレコードを`datetime`の降順でソートし、最新レコードの`type`が`checkin`かどうかで行う
- トランザクション処理が必要な場合は`db.OpenTransaction()`を使用（チェックイン時に既存チェックアウトと新規チェックインの両方を追加する場合など）
- `type`はVARCHARなので、将来的に`checkin`、`checkout`以外のタイプも追加可能
- ユーザーごとの入室記録取得では、ページネーションを実装し、チェックインレコードには対応するチェックアウト時刻も含める
