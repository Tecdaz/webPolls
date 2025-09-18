-- Tabla Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- Tabla Polls
CREATE TABLE polls (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tabla Options
CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    poll_id INTEGER NOT NULL,
    correct BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE
);

-- Índices para mejorar el rendimiento
CREATE INDEX idx_polls_user_id ON polls(user_id);
CREATE INDEX idx_options_poll_id ON options(poll_id);