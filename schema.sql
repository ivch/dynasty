create table if not exists user_roles
(
    id     serial
        constraint user_roles_pk
            primary key,
    name   varchar,
    parent int
);

insert into user_roles (id, name, parent)
values (1, 'admin', 0),
       (2, 'service', 1),
       (3, 'guard', 2),
       (4, 'neighbor', 1);

create table if not exists buildings
(
    id      serial
        constraint buildings_pk
            unique,
    name    varchar,
    address varchar
);
insert into buildings (name, address)
values ('Дом 1', 'Липковского 37-Г'),
       ('Дом 2', 'Липковского 37-Б');

create table users
(
    id          serial not null
        constraint user_pk
            unique,
    apartment   integer,
    email       varchar,
    password    varchar,
    phone       varchar,
    first_name  varchar,
    last_name   varchar,
    role        integer    default 4
        constraint user_user_roles_id_fk
            references user_roles
            on update cascade on delete set null,
    building_id integer
        constraint user_buildings_id_fk
            references buildings (id),
    parent_id   integer
        constraint users_users_id_fk
            references users (id)
            on update cascade on delete cascade,
    active      bool       default true,
    reg_code    varchar(5) default null
);

create index user_phone_index
    on users (phone);

create table sessions
(
    id            serial  not null
        constraint sessions_pk
            primary key,
    user_id       integer not null
        constraint sessions_users_id_fk
            references users (id)
            on update cascade on delete cascade,
    refresh_token uuid    not null,
    expires_in    bigint,
    created_at    timestamp default now(),
    updated_at    timestamp default now()
);

alter table sessions
    owner to postgres;

create index sessions_refresh_token_index
    on sessions (refresh_token);

create table reg_codes
(
    id   serial,
    code varchar(10) not null,
    used bool default false
);

create unique index reg_codes_code_uindex
    on reg_codes (code);

create index reg_codes_code_used_index
    on reg_codes (code, used);

create table requests
(
    id          serial
        constraint requests_pk
            primary key,
    type        varchar(50)               not null,
    user_id     int                       not null
        constraint requests_users_id_fk
            references users (id)
            on delete cascade,
    time        int                       not null,
    description varchar(1000),
    status      varchar(15) default 'new' not null
);

alter table users
    add active bool default true;