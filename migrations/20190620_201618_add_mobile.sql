-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE mobiles
(
    id         BIGSERIAL PRIMARY KEY,
    server_id  BIGINT
        CONSTRAINT server_2_servers_fk REFERENCES servers (id) on update cascade on delete cascade,
    token      UUID        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL

);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS mobiles CASCADE;

