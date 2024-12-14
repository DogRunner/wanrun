CREATE TABLE IF NOT EXISTS dogrun_manager_credentials (
    credential_id serial primary key,             -- PK
    auth_dogrun_manager_id bigint not null,       -- auth_dogrun_managerへの外部キー
    email varchar(255) unique not null,           -- emailはユニークでNULLを許可
    password varchar(256) not null,               -- パスワード認証用のパスワード。OAuth認証の場合はNULL。
    is_admin boolean not null,                    -- adminかの識別
    login_at timestamp not null                   -- 最後のログイン時間
);
