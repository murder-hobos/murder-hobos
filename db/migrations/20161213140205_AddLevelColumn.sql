
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE CharacterLevels ADD COLUMN level CHAR(2);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE CharacterLevels DROP COLUMN level;
