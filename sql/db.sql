CREATE TABLE IF NOT EXISTS segments (
    id serial primary key,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    username   varchar(255) not null unique,
    password   varchar(255) not null
);

CREATE TABLE IF NOT EXISTS users_segments (
    user_id int not null,
    segment_name varchar(255) not null,
    foreign key (user_id) references users (id),
    foreign key (segment_name) references segments (name)
);

CREATE UNIQUE INDEX users_segment_index on users_segments (user_id, segment_name);