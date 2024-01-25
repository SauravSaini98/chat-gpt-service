-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_gpt_responses (
    id SERIAL PRIMARY KEY,
    engine VARCHAR(255),
    prompt VARCHAR(255),
    answer TEXT,
    image_url VARCHAR(255),
    success BOOLEAN,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_chat_gpt_responses_deleted_at ON chat_gpt_responses (deleted_at);
CREATE INDEX idx_chat_gpt_responses_engine ON chat_gpt_responses (engine);
CREATE INDEX idx_chat_gpt_responses_prompt ON chat_gpt_responses (prompt);
CREATE INDEX idx_chat_gpt_responses_success ON chat_gpt_responses (success);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_gpt_responses;
-- +goose StatementEnd
