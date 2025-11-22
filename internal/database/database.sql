


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(20) DEFAULT 'editor' CHECK (role IN ('owner', 'admin', 'editor')),
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, email)
);

CREATE TABLE forms (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    
    title VARCHAR(500) NOT NULL,
    description TEXT,
    current_version INTEGER DEFAULT 1,
    
    -- settings
    allow_guest BOOLEAN DEFAULT TRUE,
    require_login BOOLEAN DEFAULT FALSE,
    
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    public_url VARCHAR(500) UNIQUE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE form_versions (
    id SERIAL PRIMARY KEY,
    form_id INTEGER NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    
    title VARCHAR(500) NOT NULL,
    description TEXT,
    fields JSONB NOT NULL,
    single_submission BOOLEAN DEFAULT FALSE,
    
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(form_id, version_number)
);

CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    form_id INTEGER NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    form_version_id INTEGER NOT NULL REFERENCES form_versions(id) ON DELETE CASCADE,
    
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_guest BOOLEAN DEFAULT FALSE,
    guest_email VARCHAR(255),
    
    answers JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);