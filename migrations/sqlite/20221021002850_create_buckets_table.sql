-- +goose Up
-- +goose StatementBegin
CREATE TABLE buckets
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    limit_code INTEGER NOT NULL
        CONSTRAINT buckets_limits_code_fk
            REFERENCES limits
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    value      TEXT    not null,
    count      INTEGER NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE buckets;
-- +goose StatementEnd
