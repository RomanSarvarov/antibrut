-- +goose Up
-- +goose StatementBegin
CREATE TABLE ip_rules
(
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    type   INTEGER NOT NULL,
    subnet TEXT    NOT NULL
);

CREATE UNIQUE INDEX ip_rules_subnet_uindex
    ON ip_rules (subnet);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ip_rules;
-- +goose StatementEnd
