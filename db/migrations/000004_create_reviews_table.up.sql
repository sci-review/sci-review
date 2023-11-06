
CREATE TABLE reviews(
    id UUID,
    title VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    archived BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT reviews_pk PRIMARY KEY (id)
);

CREATE TABLE reviewers(
    id UUID,
    user_id UUID NOT NULL,
    review_id UUID NOT NULL,
    role VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT reviewers_pk PRIMARY KEY (id),
    CONSTRAINT reviewers_fk1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT reviewers_fk2 FOREIGN KEY (review_id) REFERENCES reviews(id)
);