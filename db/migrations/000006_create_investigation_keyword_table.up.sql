CREATE TABLE investigation_keywords(
    id UUID,
    user_id UUID NOT NULL,
    investigation_id UUID NOT NULL,
    word VARCHAR NOT NULL,
    synonyms VARCHAR[] NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT investigation_keywords_pk PRIMARY KEY (id),
    CONSTRAINT investigation_keywords_fk1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT investigation_keywords_fk2 FOREIGN KEY (investigation_id) REFERENCES preliminary_investigations(id)
);