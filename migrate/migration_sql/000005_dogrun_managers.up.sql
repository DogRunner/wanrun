CREATE TABLE IF NOT EXISTS dogrun_managers (
    dogrun_manager_id serial primary key,    -- PK
    organization_id bigint not null,
    name varchar(128) not null,              -- dogrun_managerの名前
    image text,                              -- dogrun_managerの写真
    sex char(1),                             -- 性別
    reg_at timestamp not null,               -- 登録日
    upd_at timestamp not null                -- 更新日
);
