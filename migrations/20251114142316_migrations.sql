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

CREATE TABLE lobbies (
    uuid UUID,
    creator_uuid UUID,
    game_id INTEGER
);

CREATE TABLE players (
    uuid UUID,
    user_name TEXT
);

CREATE TABLE player_answers (
    id SERIAL PRIMARY KEY,
    game_id INTEGER,
    player_uuid UUID,
    question_id INTEGER
);

CREATE TABLE player_results (
    id SERIAL PRIMARY KEY,
    game_id INTEGER,
    player_uuid UUID,
    question_id INTEGER,
    score INTEGER
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS questions;

-- +goose StatementEnd
