CREATE TABLE pareceres (
    id VARCHAR(255) PRIMARY KEY,
    "user" VARCHAR(255) NOT NULL,
    creci VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    content TEXT NOT NULL,
    UNIQUE (id)
);