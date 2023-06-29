create table user_verification
(
    id          integer not null unique,
    date_update date,
    faculty     smallint,
    course      smallint,
    group_id    integer,
    idcard      integer not null
        primary key,
    groupname   integer,
    studentid   integer
);

