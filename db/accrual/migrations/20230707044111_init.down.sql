begin transaction;

drop function if exists  public.accrual_check_update CASCADE;
drop trigger if exists emp_stamp ON accrual CASCADE;

drop table if exists public.accrual;

commit;
