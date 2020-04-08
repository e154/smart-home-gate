-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
alter table mobiles
    add column request_id text;

update mobiles
    set request_id='vulnfzgTVxrJrUjAmWTe';

alter table mobiles
    alter column request_id set not null;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table mobiles
    drop column request_id cascade;
