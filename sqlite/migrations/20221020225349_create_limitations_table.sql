-- +goose Up
-- +goose StatementBegin
CREATE TABLE limitations
(
    code     TEXT PRIMARY KEY,
    max_attempts INTEGER NOT NULL,
    interval_sec  INTEGER NOT NULL
);

INSERT INTO limitations (code, max_attempts, interval_sec)
VALUES ('login', 10, 60),
       ('password', 100, 60),
       ('ip', 1000, 60)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE limitations;
-- +goose StatementEnd
