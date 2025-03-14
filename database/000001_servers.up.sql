CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL
);

INSERT INTO 
    servers (name, url) 
VALUES 
    ('User Management Service', 'http://95.85.125.16:18000/public/user-management/health'),
    ('Content Communication Service', 'http://95.85.125.16:18000/public/cc/health');