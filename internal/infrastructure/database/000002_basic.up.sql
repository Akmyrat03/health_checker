CREATE TABLE IF NOT EXISTS basic_config (
    id SERIAL PRIMARY KEY,
    check_interval INT NOT NULL,
    timeout INT NOT NULL,
    error_interval INT NOT NULL
);