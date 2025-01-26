CREATE TABLE IF NOT EXISTS organizations (
    organization_id serial primary key,          -- ユニークな組織ID
    -- dogrun_manager_id bigint not null,           -- dogrunmanagerのID
    organization_name varchar(128) not null,     -- 組織名
    contact_email varchar(256) not null,         -- 問い合わせ用メールアドレス
    phone_number varchar(15) not null,           -- 問い合わせ用電話番号
    address varchar(256) not null,               -- 住所
    description varchar(512),                    -- 組織の説明や特徴
    reg_at timestamp not null,
    upd_at timestamp not null
    -- is_verified BOOLEAN DEFAULT FALSE         -- 登録が承認済みかどうか
);
