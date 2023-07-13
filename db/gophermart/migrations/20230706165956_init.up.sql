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
--     uploaded_at  timestamp not null default now(),
--     processed_at timestamp
);

-- CREATE FUNCTION gophermart_order_stamp() RETURNS trigger AS $emp_stamp$
-- BEGIN
--     -- Check that empname and salary are given
--     IF NEW.empname IS NULL THEN
--         RAISE EXCEPTION 'empname cannot be null';
--     END IF;
--     IF NEW.salary IS NULL THEN
--         RAISE EXCEPTION '% cannot have null salary', NEW.empname;
--     END IF;
--
--     -- Who works for us when they must pay for it?
--     IF NEW.salary < 0 THEN
--         RAISE EXCEPTION '% cannot have a negative salary', NEW.empname;
--     END IF;
--
--     -- Remember who changed the payroll when
--     NEW.last_date := current_timestamp;
--     NEW.last_user := current_user;
--     RETURN NEW;
-- END;
-- $emp_stamp$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER emp_stamp BEFORE INSERT OR UPDATE ON emp
--     FOR EACH ROW EXECUTE FUNCTION emp_stamp();

create index gophermart_order_username_index_hash on public.gophermart_order using hash (username);
create index gophermart_order_status_index_hash on public.gophermart_order using hash (status);

create table if not exists public.gophermart_balance
(
    id           bigserial primary key    not null,
    username     varchar(255)             not null,
    current      int,
    withdrawn    int,
    uploaded_at  timestamp with time zone not null default now(),
    processed_at timestamp with time zone
);

create index gophermart_balance_balance_username_index_hash on public.gophermart_balance using hash (username);

commit;
