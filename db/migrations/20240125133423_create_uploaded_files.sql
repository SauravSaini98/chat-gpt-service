-- +goose Up
-- +goose StatementBegin
CREATE TABLE uploaded_files (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(255) UNIQUE,
    content_data BYTEA,
    content_type VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

);
CREATE INDEX idx_uploaded_files_deleted_at ON uploaded_files (deleted_at);
CREATE INDEX idx_uploaded_files_hash ON uploaded_files (hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE uploaded_files;
-- +goose StatementEnd
