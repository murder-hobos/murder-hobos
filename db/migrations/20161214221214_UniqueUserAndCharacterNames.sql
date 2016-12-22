
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE User ADD CONSTRAINT username_unique UNIQUE (username);
ALTER TABLE `Character` ADD UNIQUE INDEX name_unique (user_id, name);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE User DROP INDEX username_unique;
ALTER TABLE `Character` DROP INDEX name_unique;
