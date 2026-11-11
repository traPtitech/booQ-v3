# 再設計案の評価

評価対象は [元案](architecture.md) と現行の `internal/domain`・`internal/usecase`。今回の成果物は設計評価と [実装計画](implementation-plan.md) の修正までとする。

このモックの目的は、Go の型と処理の流れでアーキテクチャの概念を説明すること。実行可能な試作、テスト、DB 接続、メモリ repository の実装は必要ない。以下はユーザーの回答を反映した評価であり、[決定事項と残る判断事項](design-decisions.md) と対応付けて説明する。

## 評価の結論

`Item → Ownership → Transaction` に貸出元の参照を統一し、所有者・数量・貸出方針を Ownership に集約する方向を推奨する。「どのような物品か」と「誰が所有し、どの条件で貸すか」を分けられ、本・備品の貸出モデルを共通化できる。

回答により、貸出のみの Transaction、独立した申請、数量管理、Policy の Disable が決定した。Ownership は (ItemID, OwnerKind, OwnerID) で識別する。認可は今回全許可とし、本・備品の分類と ISBN を削除する。

## 現行実装との比較

パスはリポジトリルート基準。現状の説明は Go コードを優先する。

| 現行の構成・振る舞い | 評価と再設計で説明すること |
| --- | --- |
| `internal/domain/borrows.go` の Transaction は OwnershipID、EquipmentTransaction は ItemID を参照する | 備品にも Ownership を用意すると、貸出元の参照を共通化できる。備品専用 Transaction を統合する案を示す |
| `internal/domain/item.go` の EquipmentDetail が Count・CountMax を持つ | 在庫と物品属性が混在している。総数を Ownership に移し、残数を貸出履歴から説明する案が自然。今回は数量管理を選択した（D03） |
| `internal/domain/ownership.go` は所有者を UserID、貸出可否を Rentable で表す | User / Group を参照する Owner が必要。元案の二つの Policy だけでは Rentable=false を表せないため D04 で補う |
| 現行 Transaction に申請・承認／拒否・貸出・返却が同居する | TransactionRequest を独立させるなら、申請と実績の区別が明確になる。一方で状態と参照が増えるため D02 で比較する |
| 現行の Ownership 更新・貸出承認は所有者の UserID で認可する | Group ID と操作ユーザー ID を分け、所有者の種類を区別する。今回の認可は全許可とする（D08） |
| `internal/usecase/borrows.go` は申請時に Rentable を確認せず、貸出時・承認時に残数を検査していない | 貸出可否・数量確認をどこが担当するかを usecase の流れに示す。競合制御や DB ロックの実装はこのモックに含めない |
| `internal/repository/ownership.go` の削除には貸出中の扱いが TODO として残る | 削除と履歴の扱いは D11 で判断する。今回は削除方式や保存方法を実装しない |
| Like・Comment・Tag・File は貸出モデルから分離されている | 基本構成を維持してモックにも含める。Like 等を Ownership に付け替える変更は行わない |

## 元案への補足

| 元案の要素 | 説明用モックでの提案 | 判断事項 |
| --- | --- | --- |
| Item | 物品の共通情報を表す。ISBN と本・備品の分類は削除する | D12 |
| Ownership | (ItemID, OwnerKind, OwnerID) をキーとし、数量・貸出条件を表す | D03、D09 |
| Owner | `OwnerRef{Kind, ID}` で User / Group を参照する値として表す。独立した Owner エンティティが必要なら比較する | D07 |
| User / Group | ユーザー識別子、グループ識別子、管理者一覧を最小限示す | D07、D08 |
| Policy | Immediate / ApprovalRequired / Disable を表す | D04 |
| Transaction | 成立した貸出として扱い、Ownership と借り手・数量を持つ案 | D01、D10 |
| TransactionRequest | 承認が必要な申請として Transaction から分離する案 | D02、D05、D06 |

## 申請を分ける場合の比較

| 選択肢 | 利点 | 説明すべき負担 |
| --- | --- | --- |
| 現行に近い一つの Transaction | 型と状態管理が少なく、申請から返却まで一つの対象を追える | 自由貸出では使わない申請状態や、成立しなかった申請も同じ型に含む |
| TransactionRequest と Transaction を分ける | 申請と成立した貸出を区別でき、自由貸出は申請を経由せず表せる | 承認と貸出作成の対応、申請履歴と貸出履歴の参照を説明する必要がある |

回答により、後者の分離案を採用する。承認と貸出作成を一つの業務操作として扱う旨をコメントすれば概念説明には足り、UnitOfWork・mutex・ロールバック処理などの実装は不要。

## 維持するモデル

Like は Item と User の関係、Comment は Item への投稿、Tag は Item に付く名前として現行の関係を維持する。File はファイルのメタデータと保存先への抽象的な依存を示す。対応する usecase も計画に含める。

この範囲は「貸出以外だから省略する」のではなく、既存に近い形で記述し、変更する部分との境界を説明するために用意する。具体的なモデル・メソッド・作業順は [実装計画](implementation-plan.md)、選択が必要な事項は [判断事項](design-decisions.md) を参照する。
