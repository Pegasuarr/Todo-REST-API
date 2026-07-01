create table if not exists todos (
     id serial primary key,
     title varchar(255) not null,
    completed boolean default false,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);