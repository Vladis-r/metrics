Create table if not exists metrics (
    name varchar(255) not null unique, 
    type varchar(10) not null,
    value DOUBLE PRECISION,
    delta BIGINT,
    hash varchar(255)
);
