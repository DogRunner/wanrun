CREATE TABLE IF NOT EXISTS s3_file_info (
    s3_file_info_id serial primary key,             -- PK
    dog_owner_id INT NOT NULL,                      -- dog_ownersのFK
    file_id VARCHAR(64) NOT NULL,                   -- 識別用UUID
    s3_version_id VARCHAR(256),                     -- S3のバージョンI
    file_size BIGINT NOT NULL,                      -- ファイルサイズ
    s3_etag VARCHAR(64) NOT NULL,                   -- S3のETag
    s3_object_key varchar(256),                     -- S3のオブジェクトキー
    system_reg_user VARCHAR(64),                    -- 登録者情報
    reg_date timestamp,
    upd_date timestamp
);
