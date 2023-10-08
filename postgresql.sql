select * from pg_tables where schemaname='public';

create table users(
    kode int,
    name varchar(100),
    message varchar(100)
)