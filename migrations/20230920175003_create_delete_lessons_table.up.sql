CREATE TABLE public.deleted_lessons
(
    id serial NOT NULL,
    groupid bigint,
    creator bigint,
    creator_platform character(5) DEFAULT 'vk',
    lesson_id integer,
    date date,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.deleted_lessons
    OWNER to botkai;