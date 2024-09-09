BEGIN;

CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
);

CREATE TYPE tender_service_type AS ENUM (
    'CONSTRUCTION',
    'DELIVERY',
    'MANUFACTURE'
);

CREATE TYPE bid_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCELED',
    'APPROVED',
    'REJECTED'
);

CREATE TYPE bid_author_type AS ENUM (
    'ORGANIZATION',
    'USER'
);

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type tender_service_type,
    status tender_status,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id UUID REFERENCES employee(id) ON DELETE SET NULL,
    version INT DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    review_description TEXT,
    reviewed_by UUID REFERENCES employee(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
