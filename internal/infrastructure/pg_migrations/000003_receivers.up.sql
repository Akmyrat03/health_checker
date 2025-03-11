CREATE TABLE IF NOT EXISTS receivers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL
);

INSERT INTO receivers (email)
VALUES 
('ovez.hojagulyyev@gmail.com'),
('sabyrowvepa@gmail.com'),
('hudaynazarowymam@gmail.com');