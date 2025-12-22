# atode-mobile

「atode（あとで）」は「あとでやる」を**もう忘れない**ためのモバイルアプリです。

このリポジトリは、Flutter（iOS）をフロントエンド、Go をバックエンドとして開発するための土台（ブループリント）です。

## リポジトリ構成

- `backend/`
  - Go バックエンド
- `frontend/`
  - Flutter アプリ（iOS ネイティブ含む）

## コンセプト / 主要機能（計画）

- **クイック追加**
  - アプリを開いて 1 秒で「やること」を登録（タイトルのみでもOK）
- **リマインダー設定**
  - 「1時間後」「今夜」「明日」など、ざっくりした時間指定
- **プッシュ通知**
  - 指定時間に `atode: 〇〇をやる時間です` を通知（将来的に FCM を利用）
- **認証**
  - 自分のタスクだけが見られる（Firebase Auth）

## システム全体像（想定データフロー）

- Flutter
  - ユーザーがログイン → Firebase からトークン取得
  - タスク作成 → Go サーバーへ送信（ヘッダーにトークン）
- Go
  - トークン検証 → DB に保存
- PostgreSQL
  - データ永続化

---

# Backend（Go）

## 環境構築（Backend）

### 必要要件

- Go（例: `go 1.25+`）
- （将来）Docker / Docker Compose
- （将来）PostgreSQL

### セキュリティ上の注意

- 秘密情報（`.env`、証明書、鍵、Firebase 設定ファイル等）は **絶対にコミットしない**でください
  - ルートの `.gitignore` で除外しています
- まずは依存を増やさず、標準ライブラリ中心で進める前提です（サプライチェーンリスク低減）

## 起動方法（Backend）

この時点のバックエンドは、最小の HTTP サーバーです。

- エントリポイント: `backend/cmd/api/main.go`
- ヘルスチェック: `GET /healthz` → `ok`

### 起動

```bash
cd backend
go run ./cmd/api
```

デフォルトで `:8080` で待ち受けます。変更したい場合は環境変数 `ATODE_API_ADDR` を使います。

例:

```bash
ATODE_API_ADDR=127.0.0.1:8080 go run ./cmd/api
```

### 動作確認

別ターミナルで以下を実行します。

```bash
curl -i http://127.0.0.1:8080/healthz
```

---

# Frontend（Flutter / iOS）

## 環境構築（Frontend）

### 必要要件

- Flutter（stable）
- Xcode（iOS 開発用）
  - Command Line Tools 有効化
- iOS Simulator

※ CocoaPods は、iOS ネイティブ依存を入れるタイミングで必要になることがあります。

### セキュリティ上の注意

- Firebase を導入する場合、`GoogleService-Info.plist` 等は **秘密情報扱い**にしてください
  - ルート `.gitignore` で除外しています

## 起動方法（Frontend）

`frontend/` は Flutter の公式テンプレートで生成されています。

### 依存関係の取得

```bash
cd frontend
flutter pub get
```

※ `flutter pub get` は外部（pub.dev）にアクセスする場合があります。

### iOS シミュレーターで起動

```bash
cd frontend
flutter run
```

起動デバイスを指定する場合は、以下も利用できます。

```bash
flutter devices
flutter run -d <device_id>
```

---

## DB 設計（計画: PostgreSQL）

最小構成は 2 テーブルを想定します。

- `users`
  - `id` (UUID)
  - `firebase_uid` (String)
  - `email` (String)
- `tasks`
  - `id` (BigInt)
  - `user_id` (UUID)
  - `title` (String)
  - `status` (String) : 未完了/完了
  - `scheduled_at` (Timestamp)

## API 設計（計画）

- `POST /tasks` : タスク追加
- `GET /tasks` : 自分のタスク一覧取得
- `PUT /tasks/:id` : タスク更新（完了、時間変更など）
- `DELETE /tasks/:id` : タスク削除

## 開発ロードマップ（推奨順）

- Phase 1: 環境構築と土台
  - Docker Compose（Go + PostgreSQL）
  - Go API（Hello World）
  - Flutter（真っ白な画面の起動）
- Phase 2: 認証（Firebase Auth）
  - Flutter ログイン
  - Go トークン検証ミドルウェア
- Phase 3: コア機能（Tasks）
  - DB / マイグレーション
  - API（作成 + 一覧）
  - UI（一覧 + 追加）
- Phase 4: 通知
  - まずはローカル通知 → 必要に応じて FCM
