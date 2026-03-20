# パパクイズ

もうすぐ父親になる方向けに、出産・育児の基礎知識をクイズ形式で学べるWebアプリケーションです。

**URL**: https://main.d2o9xuo386rh0c.amplifyapp.com

---

## 機能

- カテゴリー別クイズ（正誤判定・解説表示）
- ユーザー登録 / ログイン / ログアウト（httpOnly Cookie + JWT認証）
- クイズ結果のDB保存
- マイページ（完了クイズ数・獲得ポイント・カテゴリー別進捗）

---

## 技術スタック

### Frontend
| 技術 | 用途 |
|---|---|
| Next.js 15 (App Router) | フレームワーク |
| TypeScript | 型安全な開発 |
| Tailwind CSS | スタイリング |
| shadcn/ui | UIコンポーネント |
| Firebase | Analytics |

### Backend
| 技術 | 用途 |
|---|---|
| Go (net/http) | APIサーバー |
| PostgreSQL | データベース |
| golang-migrate | DBマイグレーション |
| golang-jwt | JWT生成・検証 |
| bcrypt | パスワードハッシュ化 |
| Viper | 設定管理（環境変数対応） |

### Infrastructure (AWS)
| サービス | 用途 |
|---|---|
| AWS Amplify | フロントエンドのホスティング・CI/CD |
| AWS App Runner | バックエンドのコンテナ実行 |
| Amazon ECR | Dockerイメージのレジストリ |
| Amazon RDS (PostgreSQL) | マネージドDB |

---

## アーキテクチャ

```
GitHub
  ├── Frontend ──push──▶ AWS Amplify
  │                          │
  │                     HTTPS + Cookie (SameSite=None)
  │                          │
  └── Backend ──push──▶ Amazon ECR ──▶ AWS App Runner
                                              │
                                         SSL接続
                                              │
                                       Amazon RDS (PostgreSQL)
```

### 認証フロー

```
1. POST /api/auth/login
2. Backend がJWTを生成 → httpOnly + Secure + SameSite=None Cookie にセット
3. 以降のリクエストでブラウザが自動的にCookieを送信
4. GET /api/auth/me でセッション確認・ユーザー情報取得
```

---

## API一覧

| メソッド | エンドポイント | 説明 | 認証 |
|---|---|---|---|
| GET | `/api/quiz/:category` | カテゴリー別クイズ取得 | 不要 |
| POST | `/api/auth/signup` | 新規登録 | 不要 |
| POST | `/api/auth/login` | ログイン | 不要 |
| GET | `/api/auth/me` | ログイン中ユーザー取得 | 必要 |
| POST | `/api/auth/logout` | ログアウト | 必要 |
| POST | `/api/quiz/results` | クイズ結果保存 | 必要 |
| GET | `/api/quiz/stats` | ユーザー統計取得 | 必要 |

---

## DBスキーマ

```sql
-- クイズ問題
CREATE TABLE questions (
    id             SERIAL PRIMARY KEY,
    category       VARCHAR(100) NOT NULL,
    question       TEXT NOT NULL,
    options        JSONB NOT NULL,
    correct_answer INTEGER NOT NULL,
    explanation    TEXT NOT NULL
);

-- ユーザー
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- クイズ結果
CREATE TABLE quiz_results (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category   VARCHAR(100) NOT NULL,
    score      INTEGER NOT NULL,
    total      INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## ローカル開発環境のセットアップ

### 必要なもの
- Go 1.23+
- Node.js 20+
- Docker

### 1. リポジトリのクローン

```bash
git clone https://github.com/uruya/baby-prep-quiz.git
cd baby-prep-quiz
```

### 2. DBの起動

```bash
cd backend
docker compose up -d
```

### 3. Backendの起動

`backend/config.yaml` を作成：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: postgres
  sslmode: disable

jwt:
  secret: your-local-secret-key

app:
  frontend_url: http://localhost:3000
```

```bash
go run main.go
```

### 4. Frontendの起動

`frontend/.env.local` を作成：

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_FIREBASE_API_KEY=your-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-domain
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your-bucket
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your-sender-id
NEXT_PUBLIC_FIREBASE_APP_ID=your-app-id
```

```bash
cd frontend
npm install
npm run dev
```

---

## デプロイ

### Backend (AWS App Runner)

```bash
# ECRにログイン
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin \
  <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com

# ビルド・プッシュ
docker build -t baby-prep-quiz-backend ./backend
docker tag baby-prep-quiz-backend:latest \
  <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/baby-prep-quiz-backend:latest
docker push \
  <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com/baby-prep-quiz-backend:latest
```

App Runnerに設定する環境変数：

| キー | 説明 |
|---|---|
| `DATABASE_HOST` | RDSエンドポイント |
| `DATABASE_PORT` | `5432` |
| `DATABASE_USER` | DBユーザー名 |
| `DATABASE_PASSWORD` | DBパスワード |
| `DATABASE_DBNAME` | DB名 |
| `DATABASE_SSLMODE` | `require` |
| `JWT_SECRET` | JWT署名用シークレット |
| `APP_FRONTEND_URL` | AmplifyのURL |

### Frontend (AWS Amplify)

GitHubリポジトリと連携し、`main` ブランチへのpushで自動デプロイ。

Amplifyに設定する環境変数：

| キー | 説明 |
|---|---|
| `NEXT_PUBLIC_API_BASE_URL` | App RunnerのURL |
| `NEXT_PUBLIC_FIREBASE_*` | 各Firebase設定値 |
