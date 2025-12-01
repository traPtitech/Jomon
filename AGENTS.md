# Repository Guidelines

このガイドはリポジトリの変更に合わせて適宜更新してください。

## プロジェクト構成
- サーバー: `main.go` を起点に `router/`（Echo のエンドポイント）と `model/`（GORM モデル・DBアクセス）で構成。  
- クライアント: `client/src/` に Vue 2 + Vuetify の SPA。静的ファイルは `client/public/`。  
- インフラ/開発補助: `docker-compose.yml` と `Dockerfile`、テスト用 `server-test.yml`・モック用 `mock.yml`。CI は `.github/workflows/` に Go/Client/Image の 3 つ。  
- ドキュメントとデータ: `docs/`、一時保存は `uploads/`、永続化ボリューム用に `storage/`。

## ビルド・実行・開発コマンド
- サーバーモック+DB付き統合テスト: `make server-test`（Docker で mariadb を起動して Go テストを実行）。  
- クライアント開発: `make client` （`npm run lint` → `npm run serve`。localhost:3000 で確認）。  
- 単体ビルド: `go build ./...`（Go 1.17 推奨） / `cd client && npm run build`。  
- ローカル起動（Docker 本番イメージ確認用）: `docker-compose up -d --build`。停止は `make down`。

## コーディングスタイル
- Go: `gofmt`（タブインデント）必須。`go vet ./...` で静的検査。公開シンボルは PascalCase、テストダブルには `_test` 接尾辞。  
- Vue/TS: ESLint + Prettier 設定済み（`npm run lint`）。コンポーネントは `PascalCase.vue`、ルートや小要素は `kebab-case` ファイル名。CSS/SCSS は BEM か Vuetify のユーティリティを優先。  
- API/DB: モデル名は単数形、テーブル/カラムはスネークケース。JSON フィールドは camelCase を維持。

## テスト指針
- Go: `go test ./...` を基本。DB を使うテストは `MARIADB_USERNAME`/`PASSWORD`/`HOSTNAME`/`DATABASE` 環境変数を設定（CI と同値: `root/password/localhost:50000/jomon`）。  
- Vue: 現状は静的検証中心。UI 振る舞いを追加する際は `@vue/test-utils` + Jest を推奨し、`client/tests/unit/` に `_spec.js` で配置。  
- カバレッジはクリティカルパス（決裁フロー、認可、アップロード処理）を優先し、回帰バグは再現テストを追加。

## コミット & プルリク運用
- コミット: Git 履歴に倣い、冒頭に絵文字コードと短い現在形要約（例: `:bug: fix state bug`）。1 変更 1 コミットを基本。  
- PR: 目的と変更要約、テスト結果、関連 Issue/チケットを記載。UI 変更はスクリーンショット添付。CI（Go build/test、client lint/build）が緑になるまでドラフトを維持。  
- レビュー前チェック: `gofmt`, `go vet`, `go test ./...`, `cd client && npm run lint && npm run build` を実行し、差分が意図したものだけか確認。

## セキュリティ & 設定
- `.env`（DB資格情報など）はコミット禁止。共有が必要な場合はリポジトリ Secrets や開発用 1Password を利用。  
- アップロードファイルは `uploads/` に書き出されるため、実運用では適切なアクセス制御とウイルススキャンを設定。  
- コンテナ公開時は `IMAGE_TAG` を固定して再現可能ビルドにし、GHCR へ `gh workflow run` で明示的にトリガーする。  
- 依存更新は `go get -u` / `npm outdated && npm update` 後に CI を通し、Breaking Change の有無を CHANGELOG で確認。
