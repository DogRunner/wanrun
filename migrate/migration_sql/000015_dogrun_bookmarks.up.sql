CREATE TABLE IF NOT EXISTS dogrun_bookmarks (
    dogrun_bookmark_id serial primary key,
    dog_owner_id bigint not null,
    dogrun_id bigint not null,
    saved_at timestamp
)