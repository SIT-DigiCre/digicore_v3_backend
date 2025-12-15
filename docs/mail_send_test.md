# メール一括送信機能 テスト手順

## 概要

メール一括送信機能は、管理者権限を持つユーザーが複数のメールアドレスに対して一括でメールを送信できる機能です。SendGrid SMTPを使用してメールを送信します。

## 前提条件

### 環境変数の設定

`.env`ファイルに以下の環境変数を設定してください：

```bash
SENDGRID_API_KEY=your_sendgrid_api_key
ADMIN_EMAIL=admin@example.com
```

- `SENDGRID_API_KEY`: SendGridのAPIキー
- `ADMIN_EMAIL`: 管理者用メールアドレス（`sendToAdmin`が`true`の場合に送信先に追加される）

### 認証

- エンドポイントは管理者権限（`admin` claim）が必要です
- テスト時は、管理者権限を持つユーザーのJWTトークンを使用してください
- 開発環境では`.env`で`AUTH=disable`に設定することで認証を無効化できます

### メール設定

以下の設定はコード内で固定されています：

- 送信元アドレス: `contact@digicre.net`
- 送信元表示名: `デジクリ`
- SMTPサーバー: `smtp.sendgrid.net:465`

## エンドポイント

```
POST /mail
```

### 権限

- 管理者権限（`admin` claim）が必要

### リクエストボディ

```json
{
  "addresses": ["email1@example.com", "email2@example.com"],
  "subject": "メールタイトル",
  "body": "メール本文（テンプレート変数 {{.address}} が使用可能）",
  "sendToAdmin": false
}
```

#### フィールド説明

- `addresses` (必須): 送信先メールアドレスの配列
- `subject` (必須): メールのタイトル
- `body` (必須): メールの本文。テンプレート変数として以下が使用可能：
  - `{{.address}}`: 送信先メールアドレス
- `sendToAdmin` (オプション): `true`の場合、管理者用メールアドレス（`ADMIN_EMAIL`）にも送信。デフォルトは`false`

### レスポンス

#### 成功時 (200 OK)

```json
{
  "successCount": 2,
  "failures": []
}
```

#### 一部失敗時 (200 OK)

```json
{
  "successCount": 1,
  "failures": [
    {
      "address": "invalid@example.com",
      "error": "SMTP接続に失敗しました: dial tcp: lookup smtp.sendgrid.net: no such host"
    }
  ]
}
```

#### エラーレスポンス

- `400 Bad Request`: リクエストボディの形式が不正
- `403 Forbidden`: 管理者権限がない
- `500 Internal Server Error`: サーバー内部エラー

## テストケース

### 1. 正常系: 単一アドレスへの送信

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["test@example.com"],
    "subject": "テストメール",
    "body": "これはテストメールです。\n送信先: {{.address}}",
    "sendToAdmin": false
  }'
```

**期待されるレスポンス:**

```json
{
  "successCount": 1,
  "failures": []
}
```

**確認事項:**

- メールが正常に送信されること
- テンプレート変数（`{{.address}}`）が正しく展開されること
- 送信元が`デジクリ <contact@digicre.net>`であること

### 2. 正常系: 複数アドレスへの一括送信

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["user1@example.com", "user2@example.com", "user3@example.com"],
    "subject": "一括送信テスト",
    "body": "こんにちは、{{.address}}さん\n\nこれは一括送信のテストです。",
    "sendToAdmin": false
  }'
```

**期待されるレスポンス:**

```json
{
  "successCount": 3,
  "failures": []
}
```

**確認事項:**

- すべてのアドレスにメールが送信されること
- 各メールで`{{.address}}`が正しく展開されること

### 3. 正常系: 管理者用アドレスへの送信

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["user@example.com"],
    "subject": "管理者にも送信",
    "body": "このメールは管理者にも送信されます。",
    "sendToAdmin": true
  }'
```

**期待されるレスポンス:**

```json
{
  "successCount": 2,
  "failures": []
}
```

**確認事項:**

- `user@example.com`と`ADMIN_EMAIL`の両方にメールが送信されること

### 4. 異常系: 管理者権限がない

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer NON_ADMIN_JWT_TOKEN" \
  -d '{
    "addresses": ["test@example.com"],
    "subject": "テスト",
    "body": "テスト"
  }'
```

**期待されるレスポンス:**

```json
{
  "code": 403,
  "level": "Info",
  "message": "メール送信の権限がありません"
}
```

### 5. 異常系: 無効なメールアドレス

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["valid@example.com", "invalid-address"],
    "subject": "テスト",
    "body": "テスト"
  }'
```

**期待されるレスポンス:**

```json
{
  "successCount": 1,
  "failures": [
    {
      "address": "invalid-address",
      "error": "送信先の設定に失敗しました: ..."
    }
  ]
}
```

**確認事項:**

- 有効なアドレスには送信されること
- 無効なアドレスは`failures`配列に含まれること
- エラーメッセージが詳細に記録されること

### 6. 異常系: リクエストボディのバリデーションエラー

**リクエスト（addressesが空）:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": [],
    "subject": "テスト",
    "body": "テスト"
  }'
```

**期待されるレスポンス:**

```json
{
  "code": 400,
  "level": "Info",
  "message": "リクエストボディの解析に失敗しました。正しい形式で送信してください"
}
```

**リクエスト（必須フィールドが欠如）:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["test@example.com"]
  }'
```

**期待されるレスポンス:**

```json
{
  "code": 400,
  "level": "Info",
  "message": "リクエストボディの解析に失敗しました。正しい形式で送信してください"
}
```

### 7. 異常系: SendGrid接続エラー

**前提条件:**

- `SENDGRID_API_KEY`が無効または未設定
- ネットワーク接続の問題

**期待される動作:**

- すべてのメール送信が失敗する
- `failures`配列にすべてのアドレスとエラーメッセージが含まれる
- `successCount`が0になる

**期待されるレスポンス:**

```json
{
  "successCount": 0,
  "failures": [
    {
      "address": "user1@example.com",
      "error": "SMTP認証に失敗しました: ..."
    },
    {
      "address": "user2@example.com",
      "error": "SMTP認証に失敗しました: ..."
    }
  ]
}
```

## テンプレート変数のテスト

### テンプレート変数の展開確認

**リクエスト:**

```bash
curl -X POST http://localhost:8000/mail \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "addresses": ["test@example.com"],
    "subject": "【デジクリ】お知らせ",
    "body": "{{.address}}さんへ\n\nこのメールはテンプレート変数のテストです。",
    "sendToAdmin": false
  }'
```

**確認事項:**

- メール本文で`{{.address}}`が`test@example.com`に展開されること
