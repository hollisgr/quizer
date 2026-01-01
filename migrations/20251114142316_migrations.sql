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
        'qwe123qwe123'
);

INSERT INTO users (
    login, 
    password
    ) 
    VALUES (
        'akkerz',
        'akkerz123'
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
    uuid UUID UNIQUE,
    game_id INTEGER,
    is_started bool
);

CREATE TABLE players (
    uuid UUID UNIQUE,
    lobby_id UUID,
    user_name TEXT,
    is_admin bool
);

CREATE TABLE player_answers (
    id SERIAL PRIMARY KEY,
    lobby_uuid UUID,
    player_uuid UUID,
    question_id INTEGER,
    question_num INTEGER,
    answer_num INTEGER,
    answer_text TEXT DEFAULT 'NULL'
);

CREATE TABLE player_results (
    id SERIAL PRIMARY KEY,
    lobby_uuid UUID,
    player_uuid UUID,
    question_id INTEGER,
    question_num INTEGER,
    answer_num INTEGER,
    answer_text TEXT,
    score INTEGER
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS questions;

-- +goose StatementEnd
