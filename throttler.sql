CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists requests
(
    id       UUID        DEFAULT uuid_generate_v4() PRIMARY KEY NOT NULL,
    status   varchar(20) DEFAULT 'new'                          NOT NULL,
    response text        DEFAULT ''                             NOT NULL
);

UPDATE requests SET status = 'new' WHERE status = 'processing';