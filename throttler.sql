CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists requests (
    id UUID DEFAULT uuid_generate_v4(),
    status varchar(20),
    response text,
    PRIMARY KEY (id)
);