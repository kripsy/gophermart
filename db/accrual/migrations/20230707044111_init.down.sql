begin transaction;

drop index if exists public.accrual_status_index_hash;

drop function public.accrual_check_update CASCADE;
drop trigger if exists emp_stamp ON accrual CASCADE;

drop table if exists public.accrual;

commit;
