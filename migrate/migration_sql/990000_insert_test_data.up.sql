-- dog_owners テーブルに追加のテストデータを挿入
INSERT INTO dog_owners (name, image, sex, reg_at, upd_at) VALUES
('Emily Davis', 'https://example.com/images/emily.jpg', 'F', NOW(), NOW()),
('James Wilson', 'https://example.com/images/james.jpg', 'M', NOW(), NOW()),
('Olivia Martinez', NULL, 'F', NOW(), NOW()),
('William Taylor', 'https://example.com/images/william.jpg', 'M', NOW(), NOW());

-- auth_dog_ownersテーブルにデータを挿入
INSERT INTO auth_dog_owners (dog_owner_id, access_token, refresh_token, access_token_expiration, refresh_token_expiration, jwt_id, si_refresh_token, login_at) VALUES
(1, 'access_token_1', 'refresh_token_1', NOW() + INTERVAL '1 hour', NOW() + INTERVAL '7 days', 'session1', 'si_refresh_token_1', NOW()),
(2, NULL, NULL, NULL, NULL, 'session2', 'si_refresh_token_2', NOW()),
(3, NULL, NULL, NULL, NULL, 'session3', 'si_refresh_token_3', NOW()),
(4, 'access_token_4', 'refresh_token_4', NOW() + INTERVAL '1 hour', NOW() + INTERVAL '7 days', 'session4', 'si_refresh_token_4', NOW());

-- dog_owner_credentialsテーブルにデータを挿入
INSERT INTO dog_owner_credentials (auth_dog_owner_id, provider_name, grant_type, email, phone_number, provider_user_id, password, login_at) VALUES
(1, 'google', 'oauth', 'emily@example.com', NULL, 'google_user_1', NULL, NOW()),
(1, 'facebook', 'oauth', 'emily@example.com', NULL, 'facebook_user_2', NULL, NOW()),
(2, NULL, 'password', 'olivia@example.com', NULL, NULL, 'password_hash_2', NOW()),
(3, NULL, 'password', NULL, '0987654321', NULL, 'password_hash_3', NOW()),
(4, 'google', 'oauth', 'dev@example.com', NULL, 'google_user_4', NULL, NOW());

-- dog_type_mst テーブルにテストデータを挿入
INSERT INTO dog_type_mst (name) VALUES
('Labrador Retriever'),
('Beagle'),
('German Shepherd'),
('Bulldog');

-- dogs テーブルに追加のテストデータを挿入
INSERT INTO dogs (dog_owner_id, name, dog_type_id, weight, sex, image, reg_at, upd_at) VALUES
(1, 'Charlie', 1, 28, 'M', 'https://example.com/images/charlie.jpg', NOW(), NOW()),
(1, 'Daisy', 2, 22, 'F', 'https://example.com/images/daisy.jpg', NOW(), NOW()),
(2, 'Rocky', 3, 34, 'M', 'https://example.com/images/rocky.jpg', NOW(), NOW()),
(3, 'Sophie', 1, 30, 'F', 'https://example.com/images/sophie.jpg', NOW(), NOW()),
(4, 'Cooper', 2, 26, 'M', 'https://example.com/images/cooper.jpg', NOW(), NOW()),
(4, 'Chloe', 4, 15, 'F', 'https://example.com/images/chloe.jpg', NOW(), NOW());

-- dogruns テーブルに追加のテストデータを挿入
INSERT INTO dogruns (place_id, dogrun_manager_id, name, address, postcode, latitude, longitude, description, reg_at, upd_at) VALUES
(null, null, 'City Dog Park', '789 Dog Park Ave, Tokyo', '100-0003', 35.7000, 139.7100, 'A large park in the city for dogs.', NOW(), NOW()),
(null, null, 'Pawsome Adventure', '321 Woof St, Tokyo', '100-0004', 35.7100, 139.7200, 'Adventure park with various activities for dogs.', NOW(), NOW()),
(null, null, 'Happy Tails Park', '654 Tail Ave, Tokyo', '100-0005', 35.6800, 139.7500, 'A friendly park for dogs to socialize.', NOW(), NOW()),
('ChIJB6OcYgCHGGARNVckti3X6RE', null, 'わん!リトルガーデン', '東京都 葛飾区 何まち何丁目', '132-0022', 35.7111899, 139.8757518, 'テストデータ', now(), now());

-- regular_business_hours テーブルに追加のテストデータを挿入
INSERT INTO regular_business_hours (dogrun_id, day, open_time, close_time, is_closed, is_all_day, created_at, updated_at) VALUES
(1, 0, null, null, TRUE, FALSE, NOW(), NOW()),  -- City Dog Park
(1, 1, null, null, TRUE, FALSE, NOW(), NOW()),
(1, 2, '08:00:00', '18:00:00', FALSE, FALSE, NOW(), NOW()),
(1, 3, '08:00:00', '18:00:00', FALSE, FALSE, NOW(), NOW()),
(1, 4, '08:00:00', '18:00:00', FALSE, FALSE, NOW(), NOW()),
(1, 5, '08:00:00', '18:00:00', FALSE, FALSE, NOW(), NOW()),
(1, 6, '08:00:00', '18:00:00', FALSE, FALSE, NOW(), NOW()),

(2, 0, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),  -- Pawsome Adventure
(2, 1, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),
(2, 2, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),
(2, 3, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),
(2, 4, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),
(2, 5, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),
(2, 6, '09:00:00', '19:00:00', FALSE, FALSE, NOW(), NOW()),

(3, 0, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),  -- Happy Tails Park
(3, 1, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),
(3, 2, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),
(3, 3, null, null, TRUE, FALSE, NOW(), NOW()),
(3, 4, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),
(3, 5, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),
(3, 6, '07:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),

(4, 0, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),  -- わん!リトルガーデン
(4, 1, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),
(4, 2, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),
(4, 3, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),
(4, 4, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),
(4, 5, '09:00:00', '21:00:00', FALSE, FALSE, NOW(), NOW()),
(4, 6, null, null, FALSE, TRUE, NOW(), NOW());

-- special_business_hours テーブルに追加のテストデータを挿入
INSERT INTO special_business_hours (dogrun_id, date, open_time, close_time, is_closed, is_all_day, created_at, updated_at) VALUES
(1, '2024-09-17', '10:00:00', '15:00:00', FALSE, FALSE, NOW(), NOW()),  -- Valentine's Day
(2, '2024-10-21', '09:00:00', '17:00:00', FALSE, FALSE, NOW(), NOW()),  -- Spring Equinox
(3, '2024-10-04', '10:00:00', '14:00:00', FALSE, FALSE, NOW(), NOW()),  -- Children's Day
(3, '2024-12-24', null, null, TRUE, FALSE, NOW(), NOW()),  -- 臨時休業
(4, '2024-10-14', null, null, TRUE, FALSE, NOW(), NOW()),  -- 臨時休業
(4, '2024-10-15', null, null, TRUE, FALSE, NOW(), NOW());  -- 臨時休業

-- dogrun_images テーブルに追加のテストデータを挿入
INSERT INTO dogrun_images (dogrun_id, image, sort_order, upload_at) VALUES
(1, 'https://example.com/images/city_dog_park1.jpg', 1, NOW()),
(1, 'https://example.com/images/city_dog_park2.jpg', 2, NOW()),
(2, 'https://example.com/images/pawsome_adventure1.jpg', 1, NOW()),
(2, 'https://example.com/images/pawsome_adventure2.jpg', 2, NOW()),
(3, 'https://example.com/images/happy_tails1.jpg', 1, NOW()),
(3, 'https://example.com/images/happy_tails2.jpg', 2, NOW());

-- dogrun_tags テーブルに追加のテストデータを挿入
INSERT INTO dogrun_tags (dogrun_id, tag_id) VALUES
(1, 2),
(1, 6),
(1, 22),
(2, 5),
(2, 8),
(2, 11),
(2, 24),
(3, 4),
(3, 14),
(3, 19),
(3, 20),
(3, 25),
(4, 3),
(4, 7),
(4, 16),
(4, 19),
(4, 13),
(4, 21),
(4, 23),
(4, 26);


