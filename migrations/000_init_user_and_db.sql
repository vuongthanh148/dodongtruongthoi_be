-- Initialize user and database
-- This script runs as the default postgres user and creates the app user and database

-- Create the application user if it doesn't exist
DO $$ BEGIN
    CREATE USER dodongtruongthoi WITH PASSWORD 'secure_password_123';
EXCEPTION WHEN DUPLICATE_OBJECT THEN
    NULL;
END $$;

-- Create the database if it doesn't exist
DO $$ BEGIN
    CREATE DATABASE dodongtruongthoi OWNER dodongtruongthoi;
EXCEPTION WHEN DUPLICATE_DATABASE THEN
    NULL;
END $$;

-- Grant privileges
GRANT CREATE ON DATABASE dodongtruongthoi TO dodongtruongthoi;
ALTER DATABASE dodongtruongthoi OWNER TO dodongtruongthoi;
