-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
INSERT INTO chats (title, created_at) VALUES
('General Chat', NOW()),
('Random Chat', NOW()),
('Sports Chat', NOW());

-- Сообщения
INSERT INTO messages (chat_id, text, created_at) VALUES
(1, 'Hello everyone!', NOW()),
(1, 'Welcome to the general chat.', NOW()),
(2, 'Did you watch the movie last night?', NOW()),
(2, 'It was amazing!', NOW()),
(3, 'Who won the football match?', NOW()),
(3, 'Our team did!', NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
