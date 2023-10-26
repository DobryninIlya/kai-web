-- +goose Up
create table if not exists "Adv"
(
    id        integer not null
        primary key,
    groupid   integer,
    userid    integer,
    date      date,
    textfield text,
    ischeked  smallint default 0
);

alter table "Adv"
    owner to botkai;

create table if not exists activities
(
    id          integer not null
        primary key,
    name        char(120),
    description text,
    attachment  char(150)
);

alter table activities
    owner to botkai;

create table if not exists activity_members
(
    id                  integer not null
        primary key,
    activity_id         integer,
    rank_id             smallint,
    vk_id               integer,
    notification_events boolean,
    notification_news   boolean
);

alter table activity_members
    owner to botkai;

create table if not exists coworking_rent
(
    id      integer generated always as identity,
    tg_user bigint,
    status  smallint not null,
    date    date,
    time    time,
    primary key (id, status)
);

alter table coworking_rent
    owner to botkai;

create table if not exists g_chat
(
    id          bigint not null
        primary key,
    decor_mode  char(10),
    silent_mode boolean
);

alter table g_chat
    owner to botkai;

create table if not exists rent_status
(
    id   integer not null
        primary key,
    name char(30)
);

alter table rent_status
    owner to botkai;

create table if not exists saved_timetable
(
    groupp      integer not null
        primary key,
    date_update date,
    shedule     text
);

alter table saved_timetable
    owner to botkai;

create table if not exists saved_timetable_teacher
(
    login       char(50) not null
        primary key,
    date_update date,
    query       text
);

alter table saved_timetable_teacher
    owner to botkai;

create table if not exists task
(
    id                 integer not null
        primary key,
    groupid            integer not null,
    userid             integer not null,
    datee              date    not null,
    textfield          text    not null,
    attachment         text,
    "IsCheked"         smallint,
    user_source_intent char(270)
);

alter table task
    owner to botkai;

create table if not exists tg_users
(
    id             bigint not null
        primary key,
    name           char(50),
    lastname       char(50),
    username       char(50),
    role           smallint,
    "admLevel"     smallint,
    is_verificated boolean default false,
    groupname      integer,
    groupid        integer,
    login          char(50)
);

alter table tg_users
    owner to botkai;


create or replace trigger insert_notifyer_client_tg_trigger
    after insert
    on tg_users
    for each row
execute procedure insert_notifyer_client();

create table if not exists user_information
(
    id            integer not null
        primary key,
    role_id       smallint,
    name          char(40),
    lastname      char(40),
    fname         char(40),
    phone         char(15),
    email         char(50),
    scorecard_id  integer,
    group_id      integer,
    group_num     integer,
    date_update   date default CURRENT_DATE,
    full_describe text
);

alter table user_information
    owner to botkai;

create table if not exists users
(
    id_vk         bigint not null
        primary key,
    name          char(50),
    groupp        integer,
    distribution  smallint,
    admlevel      smallint,
    groupreal     integer  default 0,
    "dateChange"  date     default '2020-01-01'::date,
    balance       numeric  default 0,
    distr         smallint default 0,
    warn          smallint default 0,
    expiration    date,
    banhistory    smallint default 0,
    ischeked      smallint default 0,
    role          smallint,
    login         char(50),
    potok_lecture boolean  default true,
    has_own_shed  boolean  default false,
    affiliate     boolean  default false
);

alter table users
    owner to botkai;

create or replace trigger insert_notifyer_client_trigger
    after insert
    on users
    for each row
execute procedure insert_notifyer_client();

alter table users
    disable trigger insert_notifyer_client_trigger;

create or replace trigger insert_notifyer_client_vk_trigger
    after insert
    on users
    for each row
execute procedure insert_notifyer_client();

alter table users
    disable trigger insert_notifyer_client_vk_trigger;

create or replace trigger "insert_notifyer_client_VK_trigger"
    after insert
    on users
    for each row
execute procedure insert_notifyer_client_vk();

create table if not exists user_verification
(
    id          integer not null,
    date_update date,
    faculty     smallint,
    course      smallint,
    group_id    integer,
    idcard      integer not null
        primary key,
    groupname   integer,
    studentid   integer
);

alter table user_verification
    owner to botkai;

create table if not exists schema_migrations
(
    version bigint  not null
        primary key,
    dirty   boolean not null
);

alter table schema_migrations
    owner to botkai;

create table if not exists notify_clients
(
    id              integer not null
        primary key,
    destination_id  bigint,
    schedule_change boolean,
    source          char(2) default 'vk'::bpchar
);

alter table notify_clients
    owner to botkai;

create table if not exists deleted_lessons
(
    id               serial
        primary key,
    groupid          bigint,
    creator          bigint,
    creator_platform char(5) default 'vk'::bpchar,
    lesson_id        integer,
    date             date,
    uniqstring       char(100)
);

comment on table deleted_lessons is 'uiniqstring состоит из:
тип_занятия + время + дата';

alter table deleted_lessons
    owner to botkai;

create table if not exists lessons
(
    id          serial
        primary key,
    groupid     integer,
    daynum      smallint,
    lesson_data json,
    date        date default now()
);

alter table lessons
    owner to botkai;

create table if not exists news
(
    id          serial,
    header      char(80),
    description char(250),
    body        text,
    date        date    default now(),
    tag         char(25),
    preview_url char(256),
    author      integer default 0,
    ai_correct boolean DEFAULT false
);

alter table news
    owner to botkai;

create table if not exists api_clients
(
    uid         char(35) not null
        primary key,
    device_tag  char(16),
    token       char(64),
    create_date date default now()
);

alter table api_clients
    owner to botkai;

create table if not exists news_authors
(
    id   integer,
    name char(50)
);

alter table news_authors
    owner to botkai;

create table if not exists mobile_users
(
    uid       char(35) not null
        primary key,
    name      char(120),
    faculty   char(20),
    idcard    integer,
    groupname integer
);

alter table mobile_users
    owner to botkai;

-- ALTER TABLE IF EXISTS news
--     ADD COLUMN ai_correct boolean DEFAULT false;


-- +goose Down

drop table if exists "Adv";

drop table if exists activities;

drop table if exists activity_members;

drop table if exists coworking_rent;

drop table if exists g_chat;

drop table if exists rent_status;

drop table if exists saved_timetable;

drop table if exists saved_timetable_teacher;

drop table if exists task;

drop table if exists tg_users;

drop table if exists user_information;

drop table if exists users;

drop table if exists user_verification;

drop table if exists schema_migrations;

drop table if exists notify_clients;

drop table if exists deleted_lessons;

drop table if exists lessons;

drop table if exists news;

drop table if exists api_clients;

drop table if exists news_authors;

drop table if exists mobile_users;

