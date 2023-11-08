CREATE TABLE preliminary_investigations(
    id UUID,
    user_id UUID NOT NULL,
    review_id UUID NOT NULL,
    question TEXT NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT preliminary_investigations_pk PRIMARY KEY (id),
    CONSTRAINT preliminary_investigations_fk1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT preliminary_investigations_fk2 FOREIGN KEY (review_id) REFERENCES reviews(id)
);