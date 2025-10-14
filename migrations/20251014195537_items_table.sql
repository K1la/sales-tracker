-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id          SERIAL PRIMARY KEY,
    type        VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    amount      NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
    date        DATE NOT NULL,
    category    VARCHAR(50) NOT NULL,
    description TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_date ON items(date);
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
CREATE INDEX IF NOT EXISTS idx_items_type ON items(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
