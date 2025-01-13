CREATE TABLE IF NOT EXISTS dogrun_checkin (
    dogrun_checkin_id serial primary key,
    dog_id bigint not null,
    dogrun_id bigint not null,
    checkin_at timestamp,
    re_checkin_at timestamp
);

CREATE INDEX idx_dogrun_checkin_dogid_dogrunid_checkinat
ON dogrun_checkin (dog_id, dogrun_id, checkin_at);



CREATE TABLE IF NOT EXISTS dogrun_checkout (
    dogrun_checkout_id serial primary key,
    dog_id bigint not null,
    dogrun_id bigint not null,
    checkout_at timestamp,
    re_checkout_at timestamp
);

CREATE INDEX idx_dogrun_checkout_dogid_dogrunid_checkoutat
ON dogrun_checkout (dog_id, dogrun_id, checkout_at);