SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE items;
TRUNCATE TABLE books; 
TRUNCATE TABLE comments; 
TRUNCATE TABLE equipments; 
TRUNCATE TABLE likes; 
TRUNCATE TABLE ownerships; 
TRUNCATE TABLE tags;
TRUNCATE TABLE transactions_equipment;
TRUNCATE TABLE transactions;
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO items (id, name, description, img_url) VALUES
(1, 'item-id1', 'aaa', 'url'),
(2, 'item-id2 book', 'aaa', 'url'),
(3, 'item-id3 equipment', 'aaa', 'url'),
(4, 'item-id4 book equipment', 'aaa', 'url');

INSERT INTO books (item_id, code, created_at, updated_at) VALUES
(2, 9784088725093, '2020-01-01 00:00:00', '2020-01-01 00:00:00'),
(3, 9784088725093, '2020-01-01 00:00:00', '2020-01-01 00:00:00');

INSERT INTO comments (id, item_id, user_id, comment) VALUES 
(1, 1, 's9', 'comment'),
(2, 1, 'cp20', 'comment'),
(3, 2, 'ryoha', 'comment');

INSERT INTO equipments (item_id, count, count_max) VALUES
(3, 90, 100),
(4, 90, 100);

INSERT INTO likes (item_id, user_id) VALUES 
(1, 's9'),
(2, 's9'),
(1, 'cp20'),
(1, 'takku_bobshiroshiro_titech_trap');

INSERT INTO ownerships (id, item_id, user_id, rentalable, memo) VALUES
(1, 1, "s9", true, "memo1"),
(2, 1, "cp20", true, "memo2"),
(3, 1, "s9", false, "memo3");

INSERT INTO tags (name, item_id) VALUES
('tag1', 1),
('tag2', 1),
('tag2', 2),
('tag2', 3),
('tag3', 4);

INSERT INTO transactions_equipment (id, item_id, user_id, status, purpose, return_message) VALUES
(1, 3, 'ryoha', 0, 'かりたいから', NULL),
(2, 3, 'ryoha', 1, 'かりたいから', 'かえしました'),
(3, 4, 'ryoha', 1, 'かりたいから', 'かえしました');

INSERT INTO transactions (id, ownership_id, user_id, status, purpose, message, return_message) VALUES
(1, 1, 'ryoha', 0, 'かりたいから', NULL, NULL),
(2, 1, 'ryoha', 1, 'かりたいから', 'いいよ', NULL),
(3, 1, 'ryoha', 2, 'かりたいから', 'いいよ', 'ありがとう'),
(4, 1, 'ryoha', 3, 'かりたいから', 'ごめん、いまむり', NULL);