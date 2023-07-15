begin transaction;

create table if not exists public.accrual
(
    id           bigserial primary key    not null,
    number       bigint                   not null unique,
    status       text,
    accrual      int,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create function accrual_check_update() returns trigger AS
$emp_stamp$
begin
    NEW.processed_at := current_timestamp;
    return NEW;
end;
$emp_stamp$ LANGUAGE plpgsql;

create trigger emp_stamp
    before update
    on accrual
    for each row
execute function accrual_check_update();

create index accrual_status_index_hash on public.accrual using hash (status);

commit;
