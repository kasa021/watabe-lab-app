# LDAP認証の設定ガイド

## 研究室のLDAPサーバーと連携する方法

### 1. 必要な情報を収集

研究室の管理者に以下の情報を確認してください：

#### 必須情報
- **LDAPサーバーのホスト名/IPアドレス**  
  例: `ldap.example.ac.jp`, `192.168.1.100`

- **ポート番号**  
  - 通常のLDAP: 389
  - LDAPS (SSL/TLS): 636

- **Base DN** (ベースDN)  
  例: `dc=example,dc=ac,dc=jp`

- **ユーザー検索属性**  
  例: `uid` (通常はこれ), `cn`, `sAMAccountName` (Active Directoryの場合)

#### オプション情報
- **Bind DN** (検索用の管理者アカウント)  
  例: `cn=admin,dc=example,dc=ac,dc=jp`
  
- **Bind Password** (管理者アカウントのパスワード)

- **TLS/SSL証明書** (LDAFPSを使用する場合)

### 2. 環境変数の設定

`.env`ファイルに以下を設定：

```bash
# LDAP設定
LDAP_HOST=ldap.example.ac.jp
LDAP_PORT=389              # または 636 (LDAPS)
LDAP_BASE_DN=dc=example,dc=ac,dc=jp
LDAP_BIND_USER=cn=admin,dc=example,dc=ac,dc=jp
LDAP_BIND_PASS=your-password
```

### 3. 接続テスト方法

#### 方法1: ldapsearchコマンドでテスト（推奨）

```bash
# 基本的な接続テスト
ldapsearch -x -H ldap://ldap.example.ac.jp:389 \
  -b "dc=example,dc=ac,dc=jp" \
  -D "cn=admin,dc=example,dc=ac,dc=jp" \
  -w "your-password" \
  "(uid=username)"

# LDAFPSの場合
ldapsearch -x -H ldaps://ldap.example.ac.jp:636 \
  -b "dc=example,dc=ac,dc=jp" \
  -D "cn=admin,dc=example,dc=ac,dc=jp" \
  -w "your-password" \
  "(uid=username)"
```

#### 方法2: curlコマンドでAPIをテスト

```bash
# ログインAPIをテスト
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-username",
    "password": "your-password"
  }'

# 成功すると以下のようなレスポンスが返る
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "your-username",
    "display_name": "山田太郎",
    "email": "yamada@example.ac.jp",
    "role": "student"
  },
  "expires_at": "2024-12-28T10:00:00Z"
}
```

### 4. よくある問題と解決方法

#### エラー: "LDAP接続エラー: connection refused"

**原因:**
- LDAPサーバーのホスト名/IPアドレスが間違っている
- ポート番号が間違っている
- ファイアウォールでブロックされている

**解決方法:**
```bash
# 接続可能か確認
telnet ldap.example.ac.jp 389

# pingで到達可能か確認
ping ldap.example.ac.jp
```

#### エラー: "ユーザー検索エラー"

**原因:**
- Base DNが間違っている
- ユーザー検索属性が間違っている（`uid`ではなく`cn`など）

**解決方法:**
`internal/service/auth_service.go`の検索フィルタを変更：

```go
// uidの代わりにcnで検索
fmt.Sprintf("(cn=%s)", ldap.EscapeFilter(username))

// Active Directoryの場合
fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(username))
```

#### エラー: "認証失敗"

**原因:**
- パスワードが間違っている
- ユーザーのDNが正しく取得できていない

**解決方法:**
ログを確認してユーザーのDNが正しいか確認：
```go
// auth_service.goに一時的にログを追加
log.Printf("Found user DN: %s", userDN)
```

#### エラー: "TLS certificate error"

**原因:**
- 自己署名証明書を使用している
- 証明書の検証が厳格すぎる

**解決方法（開発環境のみ）:**
`internal/service/auth_service.go`の`InsecureSkipVerify`を`true`に変更：

```go
tlsConfig := &tls.Config{
    InsecureSkipVerify: true,  // 開発環境のみ
    ServerName:         s.config.LDAP.Host,
}
```

**注意:** 本番環境では必ず証明書を適切に設定してください！

### 5. Active Directoryとの連携

研究室がMicrosoft Active Directoryを使用している場合：

#### 設定例
```bash
LDAP_HOST=ad.example.ac.jp
LDAP_PORT=389
LDAP_BASE_DN=dc=example,dc=ac,dc=jp
LDAP_BIND_USER=CN=LDAP Admin,CN=Users,dc=example,dc=ac,dc=jp
LDAP_BIND_PASS=your-password
```

#### コード変更
`auth_service.go`の検索フィルタを変更：

```go
// sAMAccountNameで検索（Active Directory）
fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(username))

// または userPrincipalName
fmt.Sprintf("(userPrincipalName=%s@example.ac.jp)", ldap.EscapeFilter(username))
```

### 6. 開発環境でのテスト用LDAPサーバー

研究室のLDAPサーバーにアクセスできない場合、Dockerで
テスト用LDAPサーバーを起動できます：

```bash
# テスト用のOpenLDAPサーバーを起動
docker run -d -p 389:389 \
  -e LDAP_ORGANISATION="Test Lab" \
  -e LDAP_DOMAIN="example.com" \
  -e LDAP_ADMIN_PASSWORD="admin" \
  osixia/openldap:latest

# 環境変数を設定
LDAP_HOST=localhost
LDAP_PORT=389
LDAP_BASE_DN=dc=example,dc=com
LDAP_BIND_USER=cn=admin,dc=example,dc=com
LDAP_BIND_PASS=admin
```

### 7. セキュリティのベストプラクティス

#### 本番環境での必須設定

1. **TLS/SSLを必ず使用**
   ```bash
   LDAP_PORT=636  # LDAFPSポート
   ```

2. **証明書の検証を有効化**
   ```go
   InsecureSkipVerify: false
   ```

3. **管理者パスワードを環境変数で管理**
   ```bash
   # .envファイルは.gitignoreに追加
   echo ".env" >> .gitignore
   ```

4. **ログに機密情報を出力しない**
   ```go
   // パスワードをログに出力しない
   log.Printf("Authenticating user: %s", username)  // OK
   log.Printf("Password: %s", password)  // NG!
   ```

### 8. トラブルシューティングチェックリスト

- [ ] LDAPサーバーのホスト名/IPアドレスは正しいか？
- [ ] ポート番号は正しいか？（389 or 636）
- [ ] Base DNは正しいか？
- [ ] ユーザー検索属性は正しいか？（uid, cn, sAMAccountNameなど）
- [ ] Bind DN/Passwordは正しいか？
- [ ] ファイアウォールでブロックされていないか？
- [ ] TLS/SSL証明書は有効か？
- [ ] ldapsearchコマンドで接続テストをしたか？

### 9. 参考コマンド

#### LDAPツリー構造の確認
```bash
ldapsearch -x -H ldap://ldap.example.ac.jp:389 \
  -b "dc=example,dc=ac,dc=jp" \
  -D "cn=admin,dc=example,dc=ac,dc=jp" \
  -w "your-password" \
  -s sub "(objectClass=*)"
```

#### 特定ユーザーの情報を取得
```bash
ldapsearch -x -H ldap://ldap.example.ac.jp:389 \
  -b "dc=example,dc=ac,dc=jp" \
  -D "cn=admin,dc=example,dc=ac,dc=jp" \
  -w "your-password" \
  "(uid=username)" \
  cn mail displayName
```

---

困ったことがあれば、研究室の管理者に相談してください！

