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
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type tender_service_type,
    status tender_status,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    version INT DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status,
    tender_id INT REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id INT REFERENCES employee(id) ON DELETE SET NULL,
    version INT DEFAULT 1 CHECK (version > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_review (
    id SERIAL PRIMARY KEY,
    bid_id INT REFERENCES bid(id) ON DELETE CASCADE,
    review_description TEXT,
    reviewed_by INT REFERENCES employee(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
