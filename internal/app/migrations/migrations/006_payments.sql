-- +goose Up
CREATE TABLE IF NOT EXISTS public.payment_transactions
(
    uid character(32) NOT NULL,
    ext_id character(40),
    client_id character(35),
    date date DEFAULT NOW(),
    type character(12),
    ended boolean DEFAULT false,
    PRIMARY KEY (ext_id)
);

ALTER TABLE IF EXISTS public.payment_transactions
    OWNER to botkai;

-- +goose Down

DROP TABLE IF EXISTS public.payment_transactions;