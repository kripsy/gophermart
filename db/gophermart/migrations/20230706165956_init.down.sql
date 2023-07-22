begin transaction;

drop index if exists public.gophermart_order_username_index_hash;
drop index if exists public.gophermart_order_status_index_hash;
drop index if exists public.gophermart_balance_balance_username_index_hash;

drop function if exists public.gophermart_order_check_update CASCADE;
drop trigger if exists emp_stamp ON gophermart_order CASCADE;

drop table if exists public.gophermart_order;
drop table if exists public.gophermart_balance;

commit;
