-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE servers
(
    id         BIGSERIAL PRIMARY KEY,
    token      UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL

);

CREATE UNIQUE INDEX token_on_servers_unq ON servers (token);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS servers CASCADE;

