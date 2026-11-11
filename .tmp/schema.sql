-- MariaDB 向け概念スキーマ。既存 DB 用の migration ではない。
-- 文字列長・採番方法は説明用。ユーザーは別サービスが管理する。
-- user_id・borrower_id・decided_by は外部ユーザー ID。ローカルのユーザー表・FK は設けない。
-- 履歴保護のため FK は RESTRICT。所有者変更時の履歴保持は TODO(design D11)。

CREATE TABLE owner_groups (
    id VARCHAR(64) NOT NULL PRIMARY KEY
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE group_admins (
    group_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (group_id, user_id),
    FOREIGN KEY (group_id) REFERENCES owner_groups (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE items (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    img_url TEXT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE ownerships (
    item_id BIGINT NOT NULL,
    owner_kind ENUM('user', 'group') NOT NULL,
    owner_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0),
    policy ENUM('immediate', 'approval_required', 'disable') NOT NULL,
    memo TEXT NOT NULL,
    PRIMARY KEY (item_id, owner_kind, owner_id),
    FOREIGN KEY (item_id) REFERENCES items (id),
    INDEX idx_ownership_owner (owner_kind, owner_id)
    -- owner_kind=user の owner_id は外部サービスのユーザー ID。
    -- owner_kind=group は owner_groups.id。種類別の参照解決は業務側で行う。
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE transaction_requests (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    item_id BIGINT NOT NULL,
    owner_kind ENUM('user', 'group') NOT NULL,
    owner_id VARCHAR(64) NOT NULL,
    borrower_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    purpose TEXT NOT NULL,
    due_at DATETIME(6) NOT NULL,
    borrow_in_club_room BOOLEAN NOT NULL,
    status ENUM('pending', 'approved', 'rejected', 'cancelled') NOT NULL,
    created_at DATETIME(6) NOT NULL,
    decided_at DATETIME(6) NULL,
    decided_by VARCHAR(64) NULL,
    reply_message TEXT NOT NULL,
    FOREIGN KEY (item_id, owner_kind, owner_id)
        REFERENCES ownerships (item_id, owner_kind, owner_id),
    INDEX idx_request_ownership_status (item_id, owner_kind, owner_id, status),
    CHECK (
        (status = 'pending' AND decided_at IS NULL AND decided_by IS NULL)
        OR (status <> 'pending' AND decided_at IS NOT NULL AND decided_by IS NOT NULL)
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE transactions (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    item_id BIGINT NOT NULL,
    owner_kind ENUM('user', 'group') NOT NULL,
    owner_id VARCHAR(64) NOT NULL,
    borrower_id VARCHAR(64) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    purpose TEXT NOT NULL,
    due_at DATETIME(6) NOT NULL,
    borrow_in_club_room BOOLEAN NOT NULL,
    request_id BIGINT NULL, -- 自由貸出は NULL。申請経由は一申請につき一貸出。
    status ENUM('borrowed', 'returned') NOT NULL,
    checked_out_at DATETIME(6) NOT NULL,
    returned_at DATETIME(6) NULL,
    return_message TEXT NOT NULL,
    UNIQUE (request_id),
    FOREIGN KEY (request_id) REFERENCES transaction_requests (id),
    FOREIGN KEY (item_id, owner_kind, owner_id)
        REFERENCES ownerships (item_id, owner_kind, owner_id),
    INDEX idx_transaction_ownership_status (item_id, owner_kind, owner_id, status),
    CHECK (
        (status = 'borrowed' AND returned_at IS NULL)
        OR (status = 'returned' AND returned_at IS NOT NULL)
    )
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- 業務側で扱う制約（この DDL だけでは保証しない）：
-- ・貸出中数量の合計 <= Ownership.quantity。申請は残数を消費しない。
-- ・承認された申請と対応する貸出の作成を一体として扱う。
-- ・申請経由の貸出は、元申請と所有キー・借り手・数量・内容が一致する。
-- ・Policy、期限、状態遷移、貸出中・申請中の削除制限。
-- ・完了履歴の保持・廃止方式、所有者変更、既存申請がある Policy 変更。

CREATE TABLE likes (
    item_id BIGINT NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (item_id, user_id),
    FOREIGN KEY (item_id) REFERENCES items (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE comments (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    item_id BIGINT NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    text TEXT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    deleted_at DATETIME(6) NULL,
    FOREIGN KEY (item_id) REFERENCES items (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE tags (
    item_id BIGINT NOT NULL,
    name VARCHAR(191) NOT NULL,
    PRIMARY KEY (item_id, name),
    FOREIGN KEY (item_id) REFERENCES items (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE files (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(127) NOT NULL,
    created_at DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
