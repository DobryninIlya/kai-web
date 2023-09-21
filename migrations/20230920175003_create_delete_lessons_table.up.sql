-- Table: public.deleted_lessons

-- DROP TABLE IF EXISTS public.deleted_lessons;

CREATE TABLE IF NOT EXISTS public.deleted_lessons
(
    id integer NOT NULL DEFAULT nextval('deleted_lessons_id_seq'::regclass),
    groupid bigint,
    creator bigint,
    creator_platform character(5) COLLATE pg_catalog."default" DEFAULT 'vk'::bpchar,
    lesson_id integer,
    date date,
    uniqstring character(100) COLLATE pg_catalog."default",
    CONSTRAINT deleted_lessons_pkey PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.deleted_lessons
    OWNER to botkai;

COMMENT ON TABLE public.deleted_lessons
    IS 'uiniqstring состоит из:
тип_занятия + время + дата';