begin transaction;

create table if not exists public.gophermart_order
(
    id           bigint primary key       not null,
    username     varchar(255)             not null,
    number       text                     not null unique,
    status       text,
    accrual      decimal,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create index gophermart_order_username_index_hash on public.gophermart_order using hash (username);
create index gophermart_order_status_index_hash on public.gophermart_order using hash (status);

create table if not exists public.gophermart_balance
(
    id           bigint primary key,
    username     varchar(255)             not null,
    current      decimal,
    withdrawn    decimal,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create index gophermart_balance_balance_username_index_hash on public.gophermart_balance using hash (username);

commit;
