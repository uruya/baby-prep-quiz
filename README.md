# パパクイズ / baby-prep-quiz

もうすぐ父親になる方向けに、出産・育児の基礎知識をクイズ形式で学べる個人開発のWebアプリです。
Next.jsのフロントエンドとGoのバックエンドを分離し、認証、学習結果の保存、サブスクリプション連携を実装しています。

過去にAWSへデプロイした実績があります。現在の公開環境の稼働状況は未確認のため、本READMEではソースコードと構成・実装内容を中心に紹介します。

## 主な機能

- カテゴリー別クイズ、正誤判定、解説表示
- ユーザー登録・ログイン・ログアウト
- クイズ結果の保存、マイページでの進捗・ポイント表示
- 無料カテゴリーと、有効なサブスクリプションが必要なカテゴリーの出し分け
- Stripe Checkout・Customer Portalとの連携、Webhookによる購読状態の更新

課金連携は実装内容の紹介であり、有料サービスとしての運用実績や売上を示すものではありません。

## 技術スタック

| 領域 | 技術 |
|---|---|
| Frontend | Next.js 15 / App Router、React 19、TypeScript、Tailwind CSS、shadcn/ui |
| Backend | Go 1.23、net/http、database/sql・pgx、Viper |
| データ・認証 | PostgreSQL、golang-migrate、JWT、bcrypt |
| 外部サービス | Stripe、Firebase Analytics |
| テスト・CI | Go testing / httptest、GitHub Actions |
| 過去のAWSデプロイ構成 | Amplify、App Runner、ECR、RDS（PostgreSQL） |

アプリのログイン認証はGo側で実装しています。

## 現在のアプリケーション構成

```text
ブラウザ
  │ 同一オリジンの /api/* へリクエスト
  ▼
Next.js App Router
  │ Route Handler → proxyToBackend
  │ サーバー側の BACKEND_URL へ転送
  ▼
Go API（net/http）
  │ handler → usecase → repository
  ▼
PostgreSQL

Stripe ──署名付きWebhook──▶ Go API /api/billing/webhook
```

- [frontend/src/lib/proxy.ts](frontend/src/lib/proxy.ts)でCookie・リクエスト本文・クエリをバックエンドへ転送し、レスポンスのステータスとSet-Cookieをブラウザへ返します。204は本文なしで処理し、バックエンドに接続できない場合は502を返します。
- Go側は[handler](backend/handler)、[usecase](backend/usecase)、[repository](backend/repository)、[domain](backend/domain)に分けています。[main.go](backend/main.go)で依存関係とルーティングを組み立てます。
- SQLマイグレーションは[backend/migrations](backend/migrations)で管理し、API起動時に適用します。

### 認証・購読状態

パスワードをbcryptでハッシュ化し、ログイン時にJWTを発行します。CookieにはHttpOnly・Secure・SameSite=Noneを指定しています。Next.jsのプロキシがCookieを転送し、Go側で認証を行います。

有料カテゴリーへのアクセスは、Go側で認証とDB上の購読状態を確認します。StripeのWebhookは署名を検証し、Checkout完了・購読更新・購読削除のイベントに応じて状態を更新します。DB更新に失敗した場合は500を返し、Stripeによる再送の対象とします。

## コードを見る際の入口

| 確認できる内容 | 参照先 |
|---|---|
| 依存関係の組み立て、APIルーティング、DB接続 | [backend/main.go](backend/main.go) |
| HTTP処理、認証、課金・Webhook処理 | [backend/handler](backend/handler) |
| ビジネスロジックとテスト | [backend/usecase](backend/usecase) |
| PostgreSQLへのアクセス | [backend/repository](backend/repository) |
| DBスキーマと変更履歴 | [backend/migrations](backend/migrations) |
| FrontendからGoへの中継 | [frontend/src/lib/proxy.ts](frontend/src/lib/proxy.ts) |
| CIの設定 | [.github/workflows/ci.yml](.github/workflows/ci.yml) |

## API一覧

| メソッド | パス | 内容・条件 |
|---|---|---|
| GET | `/api/quiz/{category}` | 無料カテゴリーは認証不要。有料カテゴリーは認証・有効な購読が必要 |
| POST | `/api/auth/signup` | 新規登録 |
| POST | `/api/auth/login` | ログイン |
| GET | `/api/auth/me` | 認証済みユーザーの取得 |
| POST | `/api/auth/logout` | 認証Cookieの削除 |
| POST | `/api/quiz/results` | 結果保存（認証が必要） |
| GET | `/api/quiz/stats` | 統計取得（認証が必要） |
| GET | `/api/subscription/status` | 購読状態取得（認証が必要） |
| POST | `/api/billing/checkout` | Checkout開始（認証が必要） |
| POST | `/api/billing/portal` | Customer Portal開始（認証が必要） |
| POST | `/api/billing/webhook` | Stripeイベント受信（Webhook署名を検証） |

Webhookの送信先はGo APIです。Next.js側にWebhook用の中継ルートはありません。

## ローカル開発

### 必要なもの

Goは[go.mod](backend/go.mod)の指定、Node.jsは20系、Docker Composeを使用します。Firebaseの設定値は自身のプロジェクトの値を用意してください。課金動作を確認する場合はStripeのテスト用設定も必要です。

```bash
git clone https://github.com/uruya/baby-prep-quiz.git
cd baby-prep-quiz/backend
docker compose up -d
cp config.yaml.example config.yaml
```

`backend/config.yaml`のJWT秘密鍵をローカル用のランダムな値へ変更します。Stripeを使う場合はテスト用秘密鍵・Webhook署名シークレット・価格IDも設定します。設定例の値は実際の認証情報ではありません。

```bash
# backend ディレクトリで実行
go run .
```

APIは8080番ポートで起動します。別ターミナルでリポジトリの`frontend`ディレクトリへ移動し、`.env.local`を作成します。

```dotenv
BACKEND_URL=http://127.0.0.1:8080
NEXT_PUBLIC_FIREBASE_API_KEY=your-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-domain
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-bucket
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your-sender-id
NEXT_PUBLIC_FIREBASE_APP_ID=your-app-id
```

`BACKEND_URL`はNext.jsサーバーからGoへの接続先です。現在のプロキシは、以前のREADMEにあった`NEXT_PUBLIC_API_BASE_URL`を参照していません。

```bash
# frontend ディレクトリで実行
npm ci
npm run dev
```

画面は`http://localhost:3000`です。現在の認証CookieはSecure指定のため、ローカルHTTPではブラウザの扱いによってログイン状態を保持できない場合があります。認証を確認する際はHTTPSの開発環境を用意し、`app.frontend_url`もそのオリジンに合わせてください。

### Backendの環境変数

Viperにより設定ファイルの値を環境変数で上書きできます。

| 変数 | 用途 |
|---|---|
| `DATABASE_HOST` / `DATABASE_PORT` | PostgreSQL接続先 |
| `DATABASE_USER` / `DATABASE_PASSWORD` | DB認証 |
| `DATABASE_DBNAME` / `DATABASE_SSLMODE` | DB名・SSL設定 |
| `JWT_SECRET` | JWT署名用秘密鍵 |
| `APP_FRONTEND_URL` | Frontendのオリジン。CORS・課金画面からの戻り先に使用 |
| `APP_STRIPE_SECRET_KEY` | Stripe秘密鍵 |
| `APP_STRIPE_WEBHOOK_SECRET` | Webhook署名検証用シークレット |
| `APP_STRIPE_PRICE_ID` | 購読商品の価格ID |

秘密鍵・DBパスワードを含む設定ファイルはコミットしないでください。

## テストとCI

```bash
# backend ディレクトリで実行
go test ./... -v -race
```

- 認証のテストではパスワード不一致、ユーザー不在、不正なトークンなどを確認しています。
- 課金・Webhookのテストではhttptestとリポジトリのモックを使い、署名検証、購読イベントに応じた状態更新、DBエラー時の応答などを確認しています。
- これらは実際のStripe決済や実DBを通すE2Eテストではありません。
- GitHub ActionsにはGoのテスト、Frontendのlint・buildのジョブを定義しています。設定の存在と実行成功は別なので、最新の結果は[Actions](https://github.com/uruya/baby-prep-quiz/actions)で確認してください。

## AWSへのデプロイ経験

過去にFrontendをAmplify、Go APIをApp Runner、コンテナイメージをECR、PostgreSQLをRDSへ配置しました。現在の稼働状態は未確認です。

現在のコードをデプロイする際は、Next.jsサーバーが`BACKEND_URL`を参照できるよう設定する必要があります。GoコンテナにはDB・JWT・Stripeの環境変数を渡し、StripeのWebhook送信先には公開されたGo APIのURLを設定します。

## 今後の改善候補

- 実DB・Stripeテスト環境を含むE2Eの動作検証
- クイズ結果を保存する際の入力検証の強化
- デプロイ手順と公開環境の再現性・稼働確認の整備
