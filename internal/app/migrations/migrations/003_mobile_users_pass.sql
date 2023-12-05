-- +goose Up
CREATE TABLE IF NOT EXISTS public.mobile_user_password
(
    uid character(64) NOT NULL,
    login character(50),
    encrypted_password bytea,
    PRIMARY KEY (uid)
);

ALTER TABLE IF EXISTS public.mobile_user_password
    OWNER to botkai;

-- +goose Down

DROP TABLE IF EXISTS public.mobile_user_password;
