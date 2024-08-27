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
values ('37-В', 'Киев, Василия Липковского 37-В'),
       ('37-Б', 'Киев, Василия Липковского 37-Б');

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
    entry_id    int    null,
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
    type        varchar(50)                                  not null,
    user_id     integer                                      not null
        constraint requests_users_id_fk
            references users (id)
            on delete cascade,
    time        integer                                      not null,
    description varchar(1000),
    status      varchar(15) default 'new'::character varying not null,
    images      text[]      default '{}'::text[],
    history     text[]      default '{}'::text[],
    created_at  timestamp   default CURRENT_TIMESTAMP        not null,
    deleted_at  timestamp
);

alter table users
    add active bool default true;

create table entries
(
    id          serial,
    name        varchar(20),
    building_id int
        constraint entries_buildings_id_fk
            references buildings (id)
            on update cascade on delete cascade
);

insert into entries (name, building_id)
values ('Секцiя 1', 1),
       ('Секцiя 2', 1),
       ('Секцiя 3', 1),
       ('Секцiя 4', 1),
       ('Секцiя 5', 1),
       ('Секцiя 1', 2),
       ('Секцiя 2', 2),
       ('Секцiя 3', 2),
       ('Секцiя 4', 2);


create table password_recovery
(
    id serial
        constraint password_recovery_pk
            primary key,
    user_id int not null
        constraint password_recovery_users_id_fk
            references users (id)
            on update cascade on delete cascade,
    code varchar not null,
    created_at timestamp default current_timestamp,
    active boolean default true
);