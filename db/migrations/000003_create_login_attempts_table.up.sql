CREATE TABLE login_attempts (
    id UUID,
    user_id UUID NULL,
    email VARCHAR NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    CONSTRAINT login_attempts_pk PRIMARY KEY (id),
    CONSTRAINT login_attempts_fk FOREIGN KEY (user_id) REFERENCES users(id)
);