CREATE TABLE IF NOT EXISTS segments (
    id serial primary key,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    username   varchar(255) not null unique,
    password   varchar(255) not null
);

INSERT INTO users (username, password) VALUES ('test','password'),('test1','password2');

CREATE TABLE IF NOT EXISTS users_segments (
    user_id int not null,
    segment_name varchar(255) not null,
    expires_at timestamp,
    foreign key (user_id) references users (id) on delete cascade ,
    foreign key (segment_name) references segments (name) on delete cascade
);
CREATE UNIQUE INDEX IF NOT EXISTS users_segment_index on users_segments (user_id, segment_name);

CREATE EXTENSION IF NOT EXISTS pg_cron;
SELECT cron.schedule('0 0 * * *', $$DELETE FROM users_segments WHERE expires_at < now()$$) WHERE (SELECT count(*) FROM cron.job)=0;

CREATE TYPE operation_type AS ENUM ('ADD','DEL');
CREATE TABLE IF NOT EXISTS operations (
    id serial primary key,
    user_id int not null,
    segment_name varchar(255) not null,
    type operation_type not null,
    created_at timestamp not null default now()
);

-- DROP TABLE IF EXISTS segments;
-- DROP TABLE IF EXISTS users;
-- DROP TABLE IF EXISTS users_segments;
-- DROP TABLE IF EXISTS operations;