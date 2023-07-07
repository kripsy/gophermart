begin transaction;

drop index if exists public.accrual_status_index_hash;

drop table if exists public.accrual;

commit;