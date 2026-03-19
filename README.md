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

Cloudflare ダッシュボードの **My Profile > API Tokens** から発行してください。

| リソース | 権限 |
|---|---|
| AI Gateway | Account > AI Gateway > Edit |
| Access Applications | Account > Access: Apps and Policies > Edit |
| Access Policies | Account > Access: Apps and Policies > Edit |

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

### Access Policy

```bash
# 一覧
hashflare access policy list

# 詳細
hashflare access policy get <policy-id>
```

### エイリアス

- `ai-gateway` → `aig`
- `app` → `application`

```bash
hashflare aig list
hashflare access application list
```
