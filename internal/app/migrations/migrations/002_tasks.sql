-- +goose Up
CREATE TABLE IF NOT EXISTS public.app_tasks
(
    id serial NOT NULL,
    header character(150),
    body character(1500),
    create_date date DEFAULT now(),
    deadline_date date,
    author character(30),
    attachments character(40)[],
    groupname integer,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.app_tasks
    OWNER to botkai;


-- +goose Down

DROP TABLE IF EXISTS public.app_tasks;
