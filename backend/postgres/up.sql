create table orders
(
    uid                varchar not null
        constraint orders_pk
            primary key,
    track_number       varchar,
    entry              varchar,
    locale             varchar,
    internal_signature varchar,
    customer_id        varchar,
    delivery_service   varchar,
    shardkey           varchar,
    sm_id              integer,
    date_created       varchar,
    oof_shard          varchar
);

create table payment
(
    uid           varchar
        constraint payment_orders_uid_fk
            references orders,
    transaction   varchar,
    request_id    varchar default ''::character varying,
    currency      varchar,
    provider      varchar,
    amount        integer,
    payment_dt    integer,
    bank          varchar,
    delivery_cost integer,
    goods_total   integer,
    custom_fee    integer
);

create table delivery
(
    uid     varchar
        constraint delivery_orders_uid_fk
            references orders,
    phone   varchar,
    zip     varchar,
    city    varchar,
    address varchar,
    region  varchar,
    email   varchar,
    name    varchar
);

create table items
(
    uid          varchar
        constraint items_orders_uid_fk
            references orders,
    chrt_id      integer,
    track_number varchar,
    price        integer,
    rid          varchar,
    name         varchar,
    sale         integer,
    size         varchar,
    total_price  integer,
    nm_id        integer,
    brand        varchar,
    status       integer
);

