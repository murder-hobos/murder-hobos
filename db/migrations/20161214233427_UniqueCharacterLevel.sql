
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE CharacterLevels MODIFY level char(2) NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE CharacterLevels MODIFY level char(2) NULL;

