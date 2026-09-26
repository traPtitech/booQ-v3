# booQ

management tool for equipment and book rental

## 開発環境

必要なもの
- mise
- docker (+ docker compose)

インストール
```bash
mise install
```

### サーバー立ち上げ
```bash
docker compose up -d --build
```

### mise tasks

- `mise run test`: テスト実行
- `mise run test -- --cover`: カバレッジありでテスト実行 プロジェクトルート下にcover.htmlが生成されます
- `mise run gen-oapi`: OpenAPIのコード生成
- `mise run gen-mock`: gomockのコード生成

なお，テスト実行時にはDockerが使えるようにしておく必要があります（Docker Desktopを起動しておくなど）