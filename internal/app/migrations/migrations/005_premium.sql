-- +goose Up
CREATE TABLE IF NOT EXISTS public.premium_subscribe
(
    uid character(35) COLLATE pg_catalog."default" NOT NULL,
    expire_date date,
    date_of_last_payment date,
    total_payed integer,
    CONSTRAINT premium_subscribe_pkey PRIMARY KEY (uid)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.premium_subscribe
    OWNER to botkai;

-- +goose Down