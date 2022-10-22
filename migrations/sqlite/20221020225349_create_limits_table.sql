-- +goose Up
-- +goose StatementBegin
CREATE TABLE limits
(
    code     TEXT PRIMARY KEY,
    attempts INTEGER NOT NULL,
    per_sec  INTEGER NOT NULL
);

INSERT INTO limits (code, attempts, per_sec)
VALUES ('login', 10, 60),
       ('password', 100, 60),
       ('ip', 1000, 60)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE limits;
-- +goose StatementEnd
