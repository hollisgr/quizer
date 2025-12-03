-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT,
    password TEXT
);

INSERT INTO users (
    login, 
    password
    ) 
    VALUES (
        'admin',
        'admin'
);

CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    description TEXT,
    link TEXT,
    owner_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    number INTEGER,
    description TEXT,
    game_id INTEGER,
    answer INTEGER,
    answer_text TEXT DEFAULT 'NULL',
    cost INTEGER
);

CREATE TABLE lobby (
    id SERIAL PRIMARY KEY,
    creator_id INTEGER,
    game_id INTEGER,
    description TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS questions;

-- +goose StatementEnd
