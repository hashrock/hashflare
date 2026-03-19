# hashflare

Cloudflare の AI Gateway と Zero Trust Access Application を管理する CLI ツール。

## インストール

```bash
go install github.com/hashrock/hashflare@latest
```

## 設定

### 設定ファイル（推奨）

`~/.hashflare/setting.json` を作成:

```json
{
  "api_token": "your-cloudflare-api-token",
  "account_id": "your-account-id"
}
```

### 環境変数

環境変数が設定されている場合、設定ファイルより優先されます。

```bash
export CLOUDFLARE_API_TOKEN="your-token"
export CLOUDFLARE_ACCOUNT_ID="your-account-id"
```

### API トークンの権限

[Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens) からトークンを発行してください。

使用する機能に応じて、以下の権限を付与します。

#### 基本（AI Gateway / Access の管理）

| スコープ | リソース | 権限 |
|---|---|---|
| Account | Access: Policies | Edit |
| Account | Access: Apps | Edit |
| Account | AI Gateway | Edit |

#### トークン発行機能を使う場合（`aig token create` など）

上記に加えて、**User スコープ** の権限が必要です。

| スコープ | リソース | 権限 |
|---|---|---|
| User | API Tokens | Edit |

## 使い方

### AI Gateway

```bash
# 一覧
hashflare ai-gateway list

# 詳細
hashflare ai-gateway get <gateway-id>

# 作成
hashflare ai-gateway create <gateway-id> --collect-logs --cache-ttl 300

# 更新
hashflare ai-gateway update <gateway-id> --cache-ttl 600 --rate-limit 100 --rate-interval 60

# 削除
hashflare ai-gateway delete <gateway-id>

# トークン発行（AI Gateway Read 権限）
hashflare ai-gateway token create --name "my-aig-token"

# トークン一覧
hashflare ai-gateway token list

# トークン削除
hashflare ai-gateway token delete <token-id>
```

### Zero Trust Access Application

```bash
# 一覧
hashflare access app list

# 詳細
hashflare access app get <app-id>

# 作成
hashflare access app create --name my-app --domain my-app.example.com

# 更新
hashflare access app update <app-id> --name new-name --domain new.example.com

# 削除
hashflare access app delete <app-id>
```

### Access Policy（再利用可能ポリシー）

```bash
# 一覧
hashflare access policy list

# 詳細
hashflare access policy get <policy-id>
```

### Access Application Policy（アプリケーション単位のポリシー）

```bash
# 一覧
hashflare access app policy list <app-id>

# 詳細
hashflare access app policy get <app-id> <policy-id>

# 追加（メールアドレスで許可するユーザーを指定）
hashflare access app policy add <app-id> --name "Allow team" --email user@example.com --email admin@example.com

# 削除
hashflare access app policy delete <app-id> <policy-id>
```

### エイリアス

- `ai-gateway` → `aig`
- `app` → `application`

```bash
hashflare aig list
hashflare access application list
```

## 利用シナリオ

### /admin パスにアクセス制限をかける

Web アプリの管理画面 (`example.com/admin`) に、特定のメンバーだけアクセスできるよう制限する例です。

```bash
# 1. 初期設定
hashflare configure

# 2. /admin パスに対する Access Application を作成
hashflare access app create --name "Admin Area" --domain "example.com/admin"
# => アプリケーション ID が返される（例: 550e8400-e29b-...）

# 3. 許可するユーザーのポリシーを追加
hashflare access app policy add 550e8400-e29b-... \
  --name "Allow admins" \
  --email admin@example.com \
  --email dev@example.com

# 4. 設定を確認
hashflare access app policy list 550e8400-e29b-...
```

これにより、`example.com/admin` 以下にアクセスすると Cloudflare Access のログイン画面が表示され、許可されたメールアドレスのユーザーのみ通過できるようになります。
