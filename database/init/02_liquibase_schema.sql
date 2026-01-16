-- Run this as superuser or with CREATEDB/CREATE privileges
CREATE SCHEMA IF NOT EXISTS liquibase_log;
-- CREATE SCHEMA IF NOT EXISTS myapp_schema;

-- Optional: give ownership/privileges
GRANT ALL ON SCHEMA liquibase_log TO liquibase_admin;
-- GRANT ALL ON SCHEMA myapp_schema TO postgres;