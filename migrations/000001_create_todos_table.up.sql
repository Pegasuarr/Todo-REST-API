create table if not exists todos_user (
     id serial primary key,
     title varchar(255) not null,
    completed boolean default false,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);