
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE VIEW LevelsForCharactersView AS
SELECT CH.id AS "char_id", C.id, C.name, CL.level, C.base_class_id
FROM `Character` AS CH
JOIN CharacterLevels AS CL 
ON CH.id = CL.char_id
JOIN Class AS C 
ON CL.class_id = C.id;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP VIEW IF EXISTS LevelsForCharactersView;

