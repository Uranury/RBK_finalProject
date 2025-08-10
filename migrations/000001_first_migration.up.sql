CREATE TABLE IF NOT EXISTS users (
       id UUID PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password VARCHAR(255) NOT NULL,
       balance FLOAT DEFAULT 0.0,
       role VARCHAR(50) DEFAULT 'user',
       created_at TIMESTAMP NOT NULL,
       updated_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);