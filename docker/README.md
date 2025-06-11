CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY,
    city TEXT NOT NULL,
    joined_at TIMESTAMP NOT NULL DEFAULT now()
);