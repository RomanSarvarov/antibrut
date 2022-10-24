-- +goose Up
-- +goose StatementBegin
CREATE TABLE buckets
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    limitation_code INTEGER NOT NULL
        CONSTRAINT buckets_limitations_type_fk
            REFERENCES limitations
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    value      TEXT    NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX buckets_value_limitation_code_index
    ON buckets (value, limitation_code);

CREATE INDEX buckets_created_at_index
    ON buckets (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE buckets;
-- +goose StatementEnd
