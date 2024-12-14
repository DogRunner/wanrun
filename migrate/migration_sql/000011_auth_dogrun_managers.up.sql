CREATE TABLE IF NOT EXISTS auth_dogrun_managers (
    auth_dogrun_manager_id serial primary key,         -- PK
    dogrun_manager_id bigint not null,                  -- dog_managerへの外部キー
    jwt_id varchar(45),
    login_at timestamp
);
