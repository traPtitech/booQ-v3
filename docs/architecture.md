# booQ-v3 開発ガイド

## パッケージ構成

- domain
  - アプリケーションのコアとなる構造体を定義
  - データの永続化のためにrepositoryインターフェースを定義する
- usecase
  - domainの構造体を操作するロジックを定義
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

## 例
/get/itemsを実装する場合を考えます。

まず、/get/itemsではItemという構造体を使用するので、これをdomainパッケージに定義します。
ここでは、簡単のために以下のように定義します。
```go
package domain

type Item struct {
    ID          int
    Name        string
}
```

次に、このItem構造体を永続化させるためのrepositoryインターフェースをdomainパッケージに定義します。
```go
type ItemRepository interface {
    GetAllItems() ([]Item, error)
	// ... (Delete, Update, Createなど他のメソッドも定義)
}
```

usecaseパッケージにGetItemsというユースケースを定義します。
このusecaseは、ItemRepositoryにのみ依存し、実際のデータベースの実装には依存しません。
```go
package usecase

type ItemUsecase struct {
    itemRepo domain.ItemRepository
}

func NewItemUsecase(itemRepo domain.ItemRepository) *ItemUsecase {
    return &ItemUsecase{itemRepo: itemRepo}
}

func (u *ItemUsecase) GetItems() ([]domain.Item, error) {
    return u.itemRepo.GetAllItems()
}
```

最後に、handlerパッケージにGetItemsを定義します。
このhandlerは、HTTPリクエストを受け取り、usecaseを呼び出します。handlerはusecaseにのみ依存し、データベースの実装には依存しません。

実際には、このHandler構造体がoapi-codegenで生成されたServerInterfaceを実装します。
```go
package handler

type Handler struct {
    itemUsecase *usecase.ItemUsecase
}

func NewHandler(itemUsecase *usecase.ItemUsecase) *Handler {
    return &Handler{itemUsecase: itemUsecase}
}

func (h *Handler) GetItems(c echo.Context) error {
    items, err := h.itemUsecase.GetItems()
    if err != nil {
        return c.JSON(500, map[string]string{"error": "Internal Server Error"})
    }
    return c.JSON(200, items)
}
```

また、repositoryパッケージに、ItemRepositoryインターフェースの実装を定義します。
```go
package repository

type ItemRepositoryImpl struct {
    db *sqlx.DB
}

func NewItemRepositoryImpl(db *sqlx.DB) *ItemRepositoryImpl {
    return &ItemRepositoryImpl{db: db}
}

func (r *ItemRepositoryImpl) GetAllItems() ([]domain.Item, error) {
    var items []domain.Item
    err := r.db.Select(&items, "SELECT id, name FROM items")
    if err != nil {
        return nil, err
    }
    return items, nil
}
```