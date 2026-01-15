-- Create schema
CREATE SCHEMA IF NOT EXISTS gerege_base;

-- Set default search path
ALTER DATABASE gerege_db SET search_path TO gerege_base, public;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
