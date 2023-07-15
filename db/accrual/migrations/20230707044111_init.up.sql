begin transaction;

create table if not exists public.accrual
(
    id           bigint primary key       not null,
    number       text                     not null unique,
    status       text,
    accrual      decimal,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create index accrual_status_index_hash on public.accrual using hash (status);

commit;
