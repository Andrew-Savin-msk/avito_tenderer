CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
);

CREATE TYPE bid_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCELED'
);

CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TABLE tenders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    username VARCHAR(50) NOT NULL REFERENCES employee(username) ON DELETE CASCADE
);

CREATE TABLE tenders_versions (
    id BIGSERIAL PRIMARY KEY,
    tender_id UUID NOT NULL,


    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status tender_status NOT NULL,
    type service_type NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,


    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


    FOREIGN KEY (tender_id) REFERENCES tenders(id) ON DELETE CASCADE,
    CONSTRAINT unique_condition UNIQUE (tender_id, version)
);


CREATE TYPE author_type AS ENUM (
    'User',
    'Organization'
);

CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tender_id UUID NOT NULL REFERENCES tenders(id) ON DELETE CASCADE,

    author_type author_type NOT NULL,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE
);

CREATE TABLE bids_versions (
    id BIGSERIAL PRIMARY KEY,
    bid_id UUID NOT NULL,


    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status bid_status NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,


    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


    FOREIGN KEY (bid_id) REFERENCES bids(id) ON DELETE CASCADE,
    CONSTRAINT unique_bid_version UNIQUE (bid_id, version)
);

CREATE TABLE feedbacks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID NOT NULL REFERENCES bids(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,

    feedback VARCHAR(1000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
