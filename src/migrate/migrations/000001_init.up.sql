BEGIN;

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TYPE tender_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

CREATE TYPE bid_author_type AS ENUM (
    'Organization',
    'User'
);

CREATE TYPE bid_decision_type AS ENUM (
    'Approved',
    'Rejected'
);

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status tender_status DEFAULT 'Created',
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tender_information (
    id UUID REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type tender_service_type,
    version INT DEFAULT 1 CHECK (version > 0),
	PRIMARY KEY (id, version)
);

CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status bid_status DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_information (
    id UUID REFERENCES bid(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    version INT DEFAULT 1 CHECK (version > 0),
	PRIMARY KEY (id, version)
);

CREATE TABLE bid_review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bid_decision (
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    responsible_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    decision bid_decision_type NOT NULL,
    PRIMARY KEY (bid_id, responsible_id)
);

COMMIT;
