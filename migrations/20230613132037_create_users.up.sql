create table users
(
    id_vk         integer not null
        primary key,
    name          char(50),
    groupp        integer,
    distribution  smallint,
    admlevel      smallint,
    groupreal     integer  default 0,
    "dateChange"  date     default '2020-01-01'::date,
    balance       money    default 0,
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