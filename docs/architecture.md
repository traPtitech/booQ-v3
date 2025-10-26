# architecture

## パッケージ構成

- domain
  - アプリケーションのコアとなる構造体を定義
  - データの永続化のためにrepositoryインターフェースを定義する
- usecase
  - domainの構造体を操作するロジックを実装
- repository
  - データベースとの接続を実装
  - domainで定義されたrepositoryインターフェースを実装
- handler
  - リクエストを受け取り、usecaseを呼び出す
  - oapi-codegenで生成されたコードを配置
- storage

### 依存関係

TODO: mermaidにしたい

- handler -> usecase -> domain
- usecase -> repository interface
- repository implementation -> repository interface