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

また、repositoryパッケージに、ItemRepositoryインターフェースの実装を定義します。(ここではsqlxを使用していますが、実際の実装はgormです)
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

## テスト
このアーキテクチャでは各パッケージが分離されているため、独立してテストができます。
テストでは、モックを使用してパッケージ単体でテストを行います。

### usecase, handlerのテスト
モックには[gomock](https://github.com/uber-go/mock)を使用します。

#### モックの生成
`mock_generate.go`に以下のように記述します。
```go
//go:generate mockgen -source=item.go -destination=./mock/mock_item_repository.go -package=mock_domain
```
これで、`ItemRepository`インターフェースを実装した構造体を（実際にデータベースを使用することなく）テストで利用できるようになります。

#### テストを書く
モックは、例えば以下のように使用します。
```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()

repo := mock_domain.NewMockItemRepository(ctrl)
repo.EXPECT().
    GetByID(2).
    Return(nil, domain.ErrItemNotFound).
    Times(1)
```
これは、`ItemRepository`のモックの振る舞いについて「`GetByID`メソッドが引数2で呼び出された場合に、`ErrItemNotFound`エラーを返す」と定義しています。

これを利用すると、例えば
```go
item, err := usecase.GetItemByID(2)
assert.Nil(t, item)
assert.Equal(t, domain.ErrItemNotFound, err)
```
のように実際のデータベースを使わずにテストできます。

### repositoryのテスト
repositoryのテストでは、[testcontainers](https://golang.testcontainers.org/)を使用します。

`testcontainers`の初期化は`db_test.go`で行われているので、その中にある`setupTestDB`関数を呼び出すことでテスト用のデータベースを利用できます。