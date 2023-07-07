begin transaction;

create table if not exists public.gophermart_order
(
    id           bigint primary key       not null,
    user_id      text                     not null,
    number       text                     not null unique,
    status       text,
    accrual      decimal,
    uploaded_at  timestamp with time zone not null,
    processed_at timestamp with time zone
);

create index gophermart_order_user_id_index_hash on public.gophermart_order using hash (user_id);
create index gophermart_order_status_index_hash on public.gophermart_order using hash (status);

create table if not exists public.gophermart_balance
(
    id           bigint primary key,
    user_id      text                     not null,
    current      decimal,
    withdrawn    decimal,
    uploaded_at  timestamp with time zone not null,
    processed_at timestamp with time zone
);

create index gophermart_balance_balance_user_id_index_hash on public.gophermart_balance using hash (user_id);

commit;
