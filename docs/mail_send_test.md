# メール一括送信機能 テスト手順

## 概要

メール一括送信機能は、管理者権限を持つユーザーが複数のメールアドレスに対して一括でメールを送信できる機能です。

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

## エンドポイント

### メール送信

```
POST /mail
```

### 権限

- 管理者権限（`admin` claim）が必要

### リクエストボディ

```json
{
  "addresses": ["ここにあなたのアドレスを入力"],
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
      "error": ""
    }
  ]
}
```

#### エラーレスポンス

- `400 Bad Request`: リクエストボディの形式が不正
- `403 Forbidden`: 管理者権限がない
- `500 Internal Server Error`: サーバー内部エラー
