create table if not exists users (
    id int generated always as identity primary key,
    version bigint not null default 1,
    name varchar(50) not null check ( char_length(name) between 2 and 50),
    phone varchar(12) check ( phone ~ '^[0-9]+$'),
    role varchar(10) not null check (role in ('buyer', 'seller')),
    created_at timestamptz not null default now()
);

create table if not exists products (
    id int generated always as identity primary key,
    version bigint not null default 1,
    name varchar(50) not null check ( char_length(name) between 3 and 50),
    description text,
    price int not null check ( price > 0 ),
    created_at timestamptz not null default now(),

    seller_id int not null references users(id) on delete restrict
);

create table if not exists credentials (
    user_id int primary key references users(id) on delete cascade,
    login varchar(15) not null unique,
    password_hash text not null
);

create table if not exists refresh_tokens (
    id int generated always as identity primary key,
    user_id int not null references users(id) on delete cascade,
    token_hash text not null unique,
    expires_at timestamptz not null,
    created_at timestamptz not null default now(),
    revoked_at timestamptz
);
