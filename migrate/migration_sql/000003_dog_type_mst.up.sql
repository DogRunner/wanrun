CREATE TABLE IF NOT EXISTS dog_type_mst (
    dog_type_id serial primary key,
    name varchar(64) not null
);


-- マスターデータ
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (1, '秋田犬');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (2, 'ビーグル');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (3, 'ボーダー・コリー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (4, 'ブルドッグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (5, 'チワワ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (6, 'ダックスフンド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (7, 'ダルメシアン');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (8, 'ドーベルマン');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (9, 'フレンチ・ブルドッグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (10, 'ジャーマン・シェパード');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (11, 'ゴールデン・レトリーバー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (12, 'グレート・デーン');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (13, 'グレイハウンド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (14, 'ハスキー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (15, 'ジャック・ラッセル・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (16, 'ラブラドール・レトリーバー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (17, 'マルチーズ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (18, 'ミニチュア・シュナウザー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (19, 'ポメラニアン');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (20, 'プードル');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (21, 'パグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (22, 'ロットワイラー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (23, 'サモエド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (24, '柴犬');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (25, 'シー・ズー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (26, 'シベリアン・ハスキー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (27, 'スタッフォードシャー・ブル・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (28, 'ヨークシャー・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (29, 'コッカー・スパニエル');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (30, 'ボストン・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (31, 'イングリッシュ・セッター');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (32, 'アイリッシュ・セッター');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (33, 'キャバリア・キング・チャールズ・スパニエル');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (34, 'ビション・フリーゼ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (35, 'オーストラリアン・シェパード');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (36, 'バーニーズ・マウンテン・ドッグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (37, 'ボクサー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (38, 'ウェルシュ・コーギー');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (39, 'アラスカン・マラミュート');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (40, 'アメリカン・スタッフォードシャー・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (41, 'バセット・ハウンド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (42, 'ベルジアン・マリノア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (43, 'ブラッドハウンド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (44, 'チャウ・チャウ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (45, 'イングリッシュ・マスティフ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (46, 'イタリアン・グレイハウンド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (47, 'ニューファンドランド');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (48, 'オールド・イングリッシュ・シープドッグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (49, 'パピヨン');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (50, 'セント・バーナード');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (51, 'サルーキ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (52, 'スコティッシュ・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (53, 'シェットランド・シープドッグ');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (54, 'ウェスト・ハイランド・ホワイト・テリア');
INSERT INTO dog_type_mst (dog_type_id, name) VALUES (55, 'ウィペット');

-- 初期データを考慮して、シーケンスの初期値を設定
ALTER SEQUENCE dog_type_mst_dog_type_id_seq RESTART WITH 1000;
