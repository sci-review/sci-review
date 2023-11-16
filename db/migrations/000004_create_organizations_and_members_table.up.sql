CREATE TABLE organizations(
    id UUID,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    archived BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT organizations_pk PRIMARY KEY (id)
);

CREATE TABLE members(
    id UUID,
    user_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    role VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT members_pk PRIMARY KEY (id),
    CONSTRAINT members_fk1 FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT members_fk2 FOREIGN KEY (organization_id) REFERENCES organizations(id)
);