-- +goose Up
-- +goose StatementBegin
CREATE TABLE attempts
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    bucket_id INTEGER NOT NULL
        CONSTRAINT attempts_buckets_id_fk
            REFERENCES buckets
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX attempts_created_at_bucket_id_index
    ON attempts (created_at, bucket_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE attempts;
-- +goose StatementEnd
