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

drop function if exists  public.accrual_check_update CASCADE;
drop trigger if exists emp_stamp ON accrual CASCADE;

create function accrual_check_update() returns trigger AS
$emp_stamp$
begin
    NEW.processed_at := current_timestamp;
    return NEW;
end;
$emp_stamp$ LANGUAGE plpgsql;


create trigger emp_stamp
    before update
    on public.accrual
    for each row
execute function accrual_check_update();

commit;
