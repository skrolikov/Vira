CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'user',
    confirmed BOOLEAN NOT NULL DEFAULT false,
    confirm_token TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS user_logins (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    username VARCHAR(100) NOT NULL,
    ip VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    login_time TIMESTAMP NOT NULL
);