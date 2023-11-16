CREATE TABLE refresh_tokens(
   id UUID,
   user_id UUID NOT NULL,
   parent_token_id UUID NULL,
   issued_at TIMESTAMP NOT NULL,
   expires_at TIMESTAMP NOT NULL,
   active BOOLEAN NOT NULL,
   CONSTRAINT refresh_tokens_pk PRIMARY KEY (id),
   CONSTRAINT refresh_tokens_fk FOREIGN KEY (user_id) REFERENCES users(id),
   CONSTRAINT refresh_tokens_fk2 FOREIGN KEY (parent_token_id) REFERENCES refresh_tokens(id)
);