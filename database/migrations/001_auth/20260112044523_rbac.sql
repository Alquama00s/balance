-- ==========================================================
-- RBAC Database Initialization Script (PostgreSQL)
-- ==========================================================

-- Enable UUID extension (Postgres)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================================
-- Table: privilege
-- ==========================================================
CREATE TABLE privilege (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN DEFAULT TRUE
);

-- ==========================================================
-- Table: role
-- ==========================================================
CREATE TABLE role (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN DEFAULT TRUE
);

-- ==========================================================
-- Table: "user"
-- ==========================================================
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    enabled BOOLEAN DEFAULT TRUE
);

-- ==========================================================
-- Junction Table: user_role
-- ==========================================================
CREATE TABLE user_role (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE
);

-- ==========================================================
-- Junction Table: role_privilege
-- ==========================================================
CREATE TABLE role_privilege (
    role_id UUID NOT NULL,
    privilege_id UUID NOT NULL,
    PRIMARY KEY (role_id, privilege_id),
    FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE,
    FOREIGN KEY (privilege_id) REFERENCES privilege(id) ON DELETE CASCADE
);

-- ==========================================================
-- Sample data (optional)
-- ==========================================================
-- Privileges
INSERT INTO privilege (name, description)
VALUES
('CREATE_OBJECT', 'Allows creating an object'),
('READ_OBJECT', 'Allows reading an object'),
('UPDATE_OBJECT', 'Allows updating an object'),
('DELETE_OBJECT', 'Allows deleting an object');

-- Roles
INSERT INTO role (name, description)
VALUES
('ADMIN', 'Full access'),
('EDITOR', 'Can edit content'),
('VIEWER', 'Read-only access');

-- Users
INSERT INTO "user" (username, email, password_hash)
VALUES
('alice', 'alice@example.com', 'hashed_password_here'),
('bob', 'bob@example.com', 'hashed_password_here');

-- Assign roles to users
-- Example: Alice is ADMIN, Bob is VIEWER
INSERT INTO user_role (user_id, role_id)
SELECT u.id, r.id
FROM "user" u
JOIN role r ON r.name = 'ADMIN'
WHERE u.username = 'alice';

INSERT INTO user_role (user_id, role_id)
SELECT u.id, r.id
FROM "user" u
JOIN role r ON r.name = 'VIEWER'
WHERE u.username = 'bob';

-- Assign privileges to roles
-- Example: ADMIN gets all, VIEWER gets only READ_OBJECT
INSERT INTO role_privilege (role_id, privilege_id)
SELECT r.id, p.id
FROM role r
JOIN privilege p ON r.name = 'ADMIN';

INSERT INTO role_privilege (role_id, privilege_id)
SELECT r.id, p.id
FROM role r
JOIN privilege p ON r.name = 'VIEWER' AND p.name = 'READ_OBJECT';

-- ==========================================================
-- Done
-- ==========================================================
