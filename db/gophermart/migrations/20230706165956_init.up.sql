begin transaction;

create table if not exists public.gophermart_order
(
    id           bigserial primary key    not null,
    username     varchar(255)             not null,
    number       bigint                   not null unique,
    status       text,
    accrual      int                               default 0,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

drop function if exists public.gophermart_order_check_update CASCADE;
drop trigger if exists emp_stamp ON gophermart_order CASCADE;

create function gophermart_order_check_update() returns trigger AS
$emp_stamp$
begin
    NEW.processed_at := current_timestamp;
    return NEW;
end;
$emp_stamp$ LANGUAGE plpgsql;

create trigger emp_stamp
    before update
    on gophermart_order
    for each row
execute function gophermart_order_check_update();

create index gophermart_order_username_index_hash on public.gophermart_order using hash (username);
create index gophermart_order_status_index_hash on public.gophermart_order using hash (status);

create table if not exists public.gophermart_balance
(
    id           bigserial primary key    not null,
    username     varchar(255)             not null unique,
    current      int,
    withdrawn    int,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create index gophermart_balance_balance_username_index_hash on public.gophermart_balance using hash (username);

commit;
