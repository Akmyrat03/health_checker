CREATE TABLE IF NOT EXISTS receivers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    muted BOOLEAN DEFAULT FALSE
);

INSERT INTO receivers (email)
VALUES 
('akmobile.tm@gmail.com');