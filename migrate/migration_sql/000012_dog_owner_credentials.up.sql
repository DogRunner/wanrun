CREATE TABLE IF NOT EXISTS dog_owner_credentials (
    credential_id serial primary key,             -- PK
    auth_dog_owner_id bigint not null,            -- auth_dog_ownersへの外部キー
    provider_name varchar(50),                    -- OAuthプロバイダ名（例: 'google', 'facebook' など）
    grant_type varchar(20) not null,              -- 認証方式
    email varchar(255),                           -- emailはユニークでNULLを許可
    phone_number varchar(15) unique,              -- phone_numberはユニークでNULLを許可
    provider_user_id varchar(256),                -- OAuthプロバイダから提供されるユーザーID
    password varchar(256),                        -- パスワード認証用のパスワード。OAuth認証の場合はNULL。
    login_at timestamp,                           -- 最後のログイン時間

    -- CHECK 制約: 認証方式ごとに email または phone_number が必須であることをチェック
    CHECK ((provider_name = 'password' AND (email IS NOT NULL OR phone_number IS NOT NULL)) OR provider_name != 'password')
);
