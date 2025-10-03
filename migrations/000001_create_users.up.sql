CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, 
    email VARCHAR NOT NULL UNIQUE, 
    encrypted_password VARCHAR NOT NULL
    );
