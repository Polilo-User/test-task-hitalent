-- +goose Up
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chat_id_fkey
        FOREIGN KEY (chat_id)
        REFERENCES chats (id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;