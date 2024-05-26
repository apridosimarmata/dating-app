-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    duration INT NOT NULL,
    price INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS packages;
-- +goose StatementEnd
