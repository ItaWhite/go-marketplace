create table if not exists products (
    id int generated always as identity primary key,
    version bigint not null default 1,
    name varchar(50) not null check ( char_length(name) between 3 and 50),
    price int not null check ( price >0 ),
    created_at timestamptz not null default now()
);