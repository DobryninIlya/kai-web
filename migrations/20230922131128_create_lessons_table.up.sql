CREATE TABLE public.lessons
(
    id serial NOT NULL,
    groupid integer,
    daynum smallint,
    lesson_data json,
    date date DEFAULT NOW(),
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.lessons
    OWNER to botkai;