CREATE TABLE users(
    id UUID,
    name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password CHAR(60) NOT NULL,
    role VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);