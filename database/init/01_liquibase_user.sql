-- 1. Create the user (with password)
CREATE ROLE liquibase_admin WITH LOGIN PASSWORD 'afghewDFCytr@674t367';

-- 2. Make the user able to create databases (very useful for admin)
-- ALTER ROLE liquibase_admin CREATEDB;

-- 3. Make the user able to create other roles/users
ALTER ROLE liquibase_admin CREATEROLE;

-- 4. (Optional) Give the user ability to manage replication (rarely needed)
-- ALTER ROLE app_admin REPLICATION;

-- 5. Give this user full ownership + privileges on your target database
-- Replace 'myapp_db' with your actual database name
GRANT ALL PRIVILEGES ON DATABASE balance TO liquibase_admin;

-- Bonus: Make this user the default owner of new objects in the future
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO liquibase_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO liquibase_admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO liquibase_admin;
