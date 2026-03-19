## 概要

本ドキュメントでは、以下の状態を含む部員のライフサイクルと、それを支える DB / アプリケーションロジックの仕様を定義する。

- **在籍中**
- **院進**
- **転学部 / 学籍番号変更**
- **卒業**
- **退部**
- **休学**
- **再入部**
- **留年（学年補正）**

既存実装（`user_profiles.is_member` / `is_graduated`、`active_limit`、`reentries`、`grade_updates` など）を前提にしつつ、**状態とロジックを一元的に定義する**ことを目的とする。

---

## 関連テーブルとカラム

### `users`

- `id BINARY(16) PK`
- `student_number VARCHAR(8) UNIQUE`
  - 例: `aa25001`
  - 3〜4文字目（`studentNumber[2:4]`）を「入学年度 (西暦下2桁)`YY`」とみなす

### `user_profiles`

```sql
CREATE TABLE user_profiles (
    id               BINARY(16)   NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id          BINARY(16)   NOT NULL UNIQUE,
    username         VARCHAR(255) NOT NULL,
    school_grade     INT          NOT NULL,
    icon_url         VARCHAR(255) NOT NULL,
    discord_userid   VARCHAR(255) NOT NULL DEFAULT '',
    active_limit     DATE         NOT NULL,
    is_graduated     BOOLEAN      NOT NULL DEFAULT false,
    is_member        BOOLEAN      NOT NULL DEFAULT true,
    short_introduction VARCHAR(255) NOT NULL DEFAULT 'デジクリ入りました',
    introduction       TEXT         NOT NULL,
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
```

- **主に部員状態を表すフラグ**
  - `is_member`:
    - `true`: 有効アカウント（部員としてシステム利用可能）
    - `false`: 無効アカウント（「無効なアカウントです」画面に誘導）
  - `is_graduated`:
    - `true`: 卒業済フラグ
    - `false`: 在学中
- **期限管理**
  - `active_limit DATE`:
    - 部費振込・更新により延長される「有効期限」
    - `active_limit < CURRENT_DATE` かつ `is_member = true` のレコードを一括で `is_member=false` にするバッチ (`PUT /admin/inactive`) が存在

### `reentries`（再入部申請）

```sql
CREATE TABLE reentries (
    id            BINARY(16)  NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id       BINARY(16)  NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending' | 'approved' | 'rejected' を想定
    note          VARCHAR(255) NOT NULL DEFAULT '',
    checked_by    BINARY(16)          DEFAULT NULL,
    checked_at    DATETIME            DEFAULT NULL,
    created_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_reentries_user_id_created_at (user_id, created_at),
    INDEX idx_reentries_status_created_at (status, created_at)
);
```

- ユーザーごとの再入部申請履歴を管理
- 1申請1レコードで管理し、通算回数は `COUNT(*)` で算出する
- `status` によって審査状態を管理

### `grade_updates`（学年補正申請）

```sql
CREATE TABLE grade_updates (
    id          BINARY(16) NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    user_id     BINARY(16) NOT NULL,
    grade_diff  INT        NOT NULL DEFAULT -1,
    reason      TEXT       NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending' | 'approved' | 'rejected'
    approved_by BINARY(16) NULL,
    created_at  DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_grade_updates_user_status (user_id, status, created_at),
    INDEX idx_grade_updates_status (status, created_at)
);
```

- **現状実装では `grade_diff` は常に `-1` 固定**
  - 1学年分だけ増やす（＝`school_grade = school_grade + 1`）用途
- 同一ユーザーについて、承認済みレコード (`status='approved'`) の数は最大2件まで (`maxApprovedCount = 2`)

---

## 部員状態モデル

### 論理状態とフラグの対応

| 論理状態                       | is_member | is_graduated | active_limit の関係                          | 備考                                                                           |
| ------------------------------ | --------- | ------------ | -------------------------------------------- | ------------------------------------------------------------------------------ |
| 在籍中（通常部員）             | true      | false        | `>= CURRENT_DATE`                            | すべての部員向け機能にアクセス可能                                             |
| 部費期限切れ（休学・退部状態） | false     | false        | `< CURRENT_DATE`（バッチで is_member=false） | 再入部導線を表示。API は原則禁止（一部例外あり）                               |
| 卒業済部員                     | true      | true         | 任意（期限との関係は要件に応じて決める）     | 卒業ボタンにより `is_graduated=true`。アクセス権限は要検討（下記参照）。       |
| 再入部審査中                   | false     | false        | `< CURRENT_DATE` が多い想定                  | `reentries.status='pending'` が存在。`/user/me/reentry` 以外の部員機能は不可。 |
| 退部済み（再入部上限超過など） | false     | false/true   | 任意                                         | 再入部回数制限超過などで恒久的に戻れない状態。                                 |

※ 現状の実装では、**API のゲートは `is_member` のみで判定**している（`pkg/api/authenticator/login.go` / `member_gate.go`）。  
卒業生をどこまで利用可能にするかは、`is_member` の値で制御する必要がある。

---

## アカウント作成時のロジック（学年算出）

### 現状実装

- プロフィール作成/更新は `PUT /user/me`（`pkg/user/put_user_me.go`）で行う。
- DB 側の INSERT クエリ:

```sql
-- pkg/db/sql/user/insert_user_profile.sql
INSERT INTO user_profiles (
  user_id, username, school_grade, icon_url, active_limit, short_introduction, introduction
) VALUE (
  UUID_TO_BIN(/*userId*/'...'),
  /*username*/'name',
  /*schoolGrade*/1,
  /*iconUrl*/'',
  CURRENT_DATE + INTERVAL 1 MONTH,
  /*shortIntroduction*/'intro',
  ''
);
```

- `/*schoolGrade*/1` は twowaysql のプレースホルダであり、実行時にはアプリケーション側で決定した `school_grade` が差し込まれる。
  - Google 認証によるサインアップフローでは、`pkg/google_auth/post_signup_callback.go` の `createDefaultUser` 内で
    - `enterYear := int(studentNumber[2:4])`（例: `aa25001` → `24`）
    - `schoolGrade := GetSchoolYear() - 2000 - enterYear + 1`
    - さらに先頭1文字が `'m'` / `'n'` の場合にそれぞれ +4 / +6 する
      というロジックで **初期 `school_grade` を学籍番号から自動算出**している。

### 方針: 学年の表示は常に `school_grade` を使用

- 画面上に表示する学年は、**すべて `user_profiles.school_grade` を参照する**。
- 学籍番号からの学年計算はあくまで「初期値の算出」にのみ用い、その後の変更は
  - `PUT /user/me` による編集
  - 学年補正機能（`grade_updates`）
    を通じて `school_grade` に反映する。
- これにより、「表示上の学年 = `school_grade`」という関係を常に維持する。

---

## キック / バンと再入部（休学・退部の扱い）

### active_limit と無効化 (`PUT /admin/inactive`)

- 管理者エンドポイント:

```yaml
# document/paths/admin_inactive.yml
put:
  description: "Set inactive members is_member=false in bulk"
  security:
    - BearerAuth:
        - "infra"
```

- 実体クエリ:

```sql
UPDATE user_profiles
SET is_member = false
WHERE active_limit < CURRENT_DATE
  AND is_member = true;
```

- 意味:
  - 部費有効期限 (`active_limit`) が切れたユーザーを一括で **キック（= is_member=false）** する。
  - これにより、次回以降のアクセスは `Login` ミドルウェアで「無効なアカウントです」としてブロックされる。

### 無効アカウント時の挙動（休学・退部状態）

- 認証ミドルウェア (`pkg/api/authenticator/validate.go`) の概略:
  - JWT の `sub` から `user_id` を取得。
  - `sql/user/select_user_is_member_from_user_id.sql` により `is_member` を取得。
  - `is_member=false` かつ `isNonMemberAllowedPath(c.Path())` でない場合:
    - `403 Forbidden` + メッセージ「無効なアカウントです」
- `RegisterNonMemberAllowedPath` により、**非部員でも許可されるパス**を登録可能。
  - 現状: `/user/me/payment`, `/user/me/reentry` が登録済み。

したがって、**休学や退部などで一時的に活動停止中の部員は `is_member=false` とし、例外的に `/user/me/payment` `/user/me/reentry` のみ利用可能**とする運用になっている。

### 再入部申請（休学・退部からの復帰）

#### ユーザー側エンドポイント: `PUT /user/me/reentry`

- 実装: `pkg/user/put_user_me_reentry.go`
- 主なロジック:
  1. 自分のプロフィール取得 (`GetUserProfileFromUserId`)
     - `profile.IsMember == true` の場合 → 400: 「有効なアカウントでは再入部申請できません」
  2. 未処理申請チェック
     - `sql/reentry/select_pending_reentry_count_by_user_id.sql`
     - `pending_count > 0` の場合 → 400: 「未処理の再入部申請があります」
  3. 再入部回数チェック
     - `sql/reentry/select_reentry_total_count_by_user_id.sql`
     - `maxReentryCount = 2`（`pkg/user/reentry.go`）
     - `totalCount >= 2` の場合 → 400: 「再入部申請は最大2回までです」
  4. 振込情報の更新
     - `updateUserPayment`（`PUT /user/me/payment` 共通ロジック）
  5. `reentries` へのレコード INSERT（1申請1レコード）
  6. 作成したレコードを取得し、レスポンスとして返却
  7. 学籍番号宛（`<student_number>@shibaura-it.ac.jp`）に「再入部申請受付」のメール送信

#### 管理者側エンドポイント: `PUT /admin/reentry/{reentryId}`

- 実装: `pkg/admin/put_admin_reentry_reentry_id.go`
- 主なロジック:
  1. `reentryId` から対象申請を取得
  2. `status='approved'` の場合:
     - `isCurrentSchoolYearPaymentChecked` により、**今年度の部費振込報告が会計承認済みか**確認
       - 未承認の場合 → 400: 「今年度の部費振込報告が会計承認されていないため、再入部を承認できません」
  3. `sql/reentry/update_reentry_status.sql` で `status` / `note` / `checked_by` を更新（楽観ロック: `status='pending'` のものだけ更新）
  4. `status='approved'` の場合:
     - `sql/reentry/update_user_member_status.sql` により `user_profiles.is_member=true` に更新
  5. 学籍番号を取得して、承認 / 却下結果をメール通知

#### ミドルウェア側の再入部導線制御

- `Login` ミドルウェア内で、`is_member=false && c.Path()=="/user/me/reentry"` の場合に追加チェック:
  - `hasPendingReentry(userId)` で `pending_count > 0` のとき:
    - 403 + 「再入部申請の確認中です。振込案内ページをご確認ください」
  - それ以外の非部員は `/user/me/reentry` にアクセス可能

#### 仕様要約

- **休学・退部・キックされた部員**は、次の条件を満たすと再入部できる:
  - `is_member=false`
  - 未処理の再入部申請が存在しない
  - 自身の再入部申請回数（承認・却下問わず申請レコード総数）が 2 回未満
  - 部費振込報告が行われ、今年度分の支払いが会計により承認される
- **再入部回数は部員（`user_id`）ごとに最大 2 回まで**

---

## 学年補正（留年・休学によるズレの解消）

### ユーザー側エンドポイント

- `POST /user/me/grade-update`（申請作成）
  - 実装: `pkg/user/post_user_me_grade_update.go`
- `GET /user/me/grade-update`（自分の申請一覧）
  - 実装: `pkg/user/get_user_me_grade_update.go`

### ロジック

1. **承認済み回数チェック**
   - `sql/grade_update/select_approved_count_by_user_id.sql`
   - `approved_count >= 2` の場合:
     - 400: 「学年補正申請は最大2回までです」
2. **未処理申請チェック**
   - `sql/grade_update/select_pending_count_by_user_id.sql`
   - `pending_count > 0` の場合:
     - 400: 「未処理の申請が存在します。承認/却下後に再度申請してください」
3. **申請作成**
   - `grade_diff = -1` 固定（1学年分のみ）
   - `reason` はユーザー入力
4. **管理者承認時 (`PUT /admin/grade-update/{gradeUpdateId}`)**
   - `status='pending'` → `approved` or `rejected` に更新（楽観ロック）
   - `approved` の場合:
     - 対象レコードを取得し、`user_profiles.school_grade = school_grade + grade_diff` を実行
     - 現仕様では `grade_diff=-1` 固定のため、**最大2回まで「学年+1」の補正が可能**

### 仕様要約

- **休学・留年などで「初期算出された学年」と実際の学年がズレた場合**に利用する。
- 許可される操作:
  - 1申請あたり **1学年分の増加のみ**（`grade_diff=-1`）
  - **承認済み申請はユーザーごとに最大2件**（大学の規定に合わせた上限）
- UI 上は、プロフィールページなどに「学年補正申請」フォームを設置し、理由 (`reason`) の入力を求める。

---

## 卒業の扱い

### ユーザー操作: `PUT /user/me/graduated`

- 実装: `pkg/user/put_user_me_graduated.go`
- ロジック:
  1. 自分のプロフィール取得
  2. `profile.SchoolGrade < 4` の場合:
     - 400: 「4年生以上のみ卒業済みにできます」
  3. `sql/user/update_user_is_graduated.sql` で `is_graduated = true` に更新
  4. 更新後の `GET /user/me` を返却

### is_member との関係

- 現状実装:
  - 卒業時は **`is_graduated=true` にするだけで `is_member` は変更しない**。
  - その後 `active_limit` が切れた場合、`PUT /admin/inactive` によって `is_member=false` にされる（卒業生も同じ処理対象）。
- 方針案:
  - 要件として「卒業生には引き続き一定の機能を提供したい」場合:
    - `UPDATE user_profiles SET is_member = false ...` の条件から **`is_graduated = true` を除外**する必要がある。
    - 例:
      ```sql
      UPDATE user_profiles
      SET is_member = false
      WHERE active_limit < CURRENT_DATE
        AND is_member = true
        AND is_graduated = false;
      ```
  - 逆に「卒業後も部費未納なら利用不可でよい」場合は現状のままでよい。

本仕様では、**卒業フラグ (`is_graduated`) と部員フラグ (`is_member`) を独立に扱う**ことを基本とし、実際にどう運用するかは運営ポリシーに合わせて決定する。

**TODO（将来対応）**: 卒業生については、`active_limit` を超過しても `is_member=true` を維持できるようにする（例: `PUT /admin/inactive` の WHERE 句から `is_graduated=true` を除外するなど）。

---

## 院進・転学部と学籍番号変更

### ケース1: 院進で学籍番号が変わらない場合

- `users.student_number` は変わらず、**同じ `user_id` を継続利用**する。
- 対応方針:
  - `school_grade` を 5年目・6年目…と手動更新するか、学年補正機能 (`grade_updates`) を使用する。
  - `is_graduated` は卒業時にのみ `true` にする（院進時点では `false` のまま）。

### ケース2: 転学部・再入学などで学籍番号が変わる場合

- **アカウントリンク機能は導入せず**、既存アカウントの `student_number` をインフラ権限で直接更新する方針とする。
- 想定エンドポイントの例:

```yaml
put:
  tags:
    - admin
  description: "学籍番号の更新（インフラ権限）"
  security:
    - BearerAuth:
        - "infra"
  # パス例: /admin/user/{userId}/student-number
  # リクエストボディ例:
  #   { "studentNumber": "aa25001" }
```

- このエンドポイントにより、**同一の `user_id` のまま学籍番号だけを更新**する。
  - 作品やイベント参加など `user_id` ベースの関連データはすべて引き継がれる。
  - 学年表示はあくまで `school_grade` に依存するため、必要に応じて `school_grade` も別途編集・学年補正申請で調整する。

---

## エンドポイントへの影響とルール

### 1. ログイン・アクセス制御

- `Login` ミドルウェア:
  - `is_member=false` かつ 非許可パスの場合 → 403 「無効なアカウントです」
  - 現状の非部員許可パス:
    - `/user/me/payment`
    - `/user/me/reentry`
- 今後追加が想定される非部員許可パス:
  - 再入部関連ページの参照 API

### 2. 部員一覧 (`GET /user`)

- 現状:
  - `SELECT ... FROM user_profiles` で `is_graduated` / `is_member` の値は取得するが、WHERE 句によるフィルタはない。
  - 卒業生 / 非部員の扱いはフロント側で制御している想定。
- 変更方針（例）:
  - デフォルトでは「現役部員のみ」を返す:
    - `WHERE is_member = true AND is_graduated = false`
    - さらに `user_relations.child_user_id` に存在する `user_id` は除外。
  - クエリパラメータで `includeGraduated` / `includeInactive` などを追加し、必要に応じて卒業生や無効アカウントも取得できるようにする。

### 3. 部費支払い・有効期限更新

- 部費振込報告 (`PUT /user/me/payment`):
  - 振込申請時に `active_limit` を「現在日付 + 1ヶ月」に更新（仮延長）。
- 会計承認 (`PUT /payment/{paymentId}`):
  - `checked=true` の場合、`active_limit` を `GetSchoolYear()+1` 年 `05-01` に更新。
- 手動更新 (`PUT /user/me/renewal`):
  - `active_limit` を `GetYear()` 年 `06-01` に更新。

これらにより、「部費を払っている限り毎年一定のタイミングまで `active_limit` が延長される」仕組みになっている。  
`PUT /admin/inactive` はこれを前提に、「期限切れ → is_member=false → 再入部フローへ」というライフサイクルを実現する。
