begin transaction;

drop index if exists public.gophermart_order_user_id_index_hash;
drop index if exists public.gophermart_order_status_index_hash;
drop index if exists public.gophermart_balance_balance_user_id_index_hash;

drop table if exists public.gophermart_order;
drop table if exists public.gophermart_balance;

commit;