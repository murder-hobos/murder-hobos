-------------------------------------------------------------------------------
--                              Initial Schema
--
-- A note on conventions used in this database for client applications:
--  Tables have singular names (e.g. class), 
--  while views have plural names (e.g. cannon_classes)
--  
-------------------------------------------------------------------------------

DROP VIEW IF EXISTS cannon_spells, user_spells, cannon_classes, basic_classes;
DROP TABLE IF EXISTS class_spells, character_class_levels, 
    mh_character, spell, class, mh_user, cannon_text, source;

-------------------------------------------------------------------------------
--                              Source Entities
--
-- Since both cannon texts and users can be sources for spells, classes, races,
-- etc. we want to be able to model this polymorphic association. The idea here
-- is to have all attributes common to cannon_texts and users in one parent
-- "source" table, and all unique attributes in their respective subtables.
-- In this particular instance, the only attribute they have to share is an id.
-- This id can then be referenced by any weak entity in our domain reliant on 
-- a source. 
--
-------------------------------------------------------------------------------
CREATE TABLE source (
    id SERIAL PRIMARY KEY
);

CREATE TABLE cannon_text (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    abbreviation TEXT NOT NULL,
    FOREIGN KEY (id) REFERENCES source(id) ON DELETE CASCADE
);

-- mh_user because having to quote "user" every time would get annoying for
-- client applications
CREATE TABLE mh_user (
    id INTEGER PRIMARY KEY,
    username VARCHAR(60) NOT NULL UNIQUE,
    password CHAR(60) NOT NULL, -- bcrypt length
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_active TIMESTAMP NOT NULL DEFAULT NOW(),
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (id) REFERENCES source(id)
);

-------------------------------------------------------------------------------
--                              Source Triggers
--
-- We want cannon_text and mh_user to essentially share the id column of the
-- "source" table. To this end, we define triggers here to accomplish 
-- the following:
-- 
-- Inserts:
--      Add a row to "source" and use that generated id as our new id in the
--      subtable.
--
-- Deletes:
--      After we delete from the child tables, also delete the id from the
--      parent "source" table, creating a "reverse cascade" of sorts.
-------------------------------------------------------------------------------

-- Use serial id from 'source' table for the cannon_text id
CREATE OR REPLACE FUNCTION cannon_text_insert() RETURNS TRIGGER AS $$
    BEGIN
        INSERT INTO source(id) VALUES (DEFAULT) RETURNING id INTO NEW.id;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cannon_text_insert
BEFORE INSERT ON cannon_text
    FOR EACH ROW EXECUTE PROCEDURE cannon_text_insert();

-- Delete the row from 'source' after we delete from cannon_text
CREATE OR REPLACE FUNCTION cannon_text_delete() RETURNS TRIGGER AS $$
    BEGIN
        DELETE FROM source WHERE id = OLD.id;
        IF NOT FOUND THEN RETURN NULL; END IF;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cannon_text_delete
AFTER DELETE ON cannon_text
    FOR EACH ROW EXECUTE PROCEDURE cannon_text_delete();


-- Use serial id from 'source' table for the mh_user id
CREATE OR REPLACE FUNCTION mh_user_insert() RETURNS TRIGGER AS $$
    BEGIN
        INSERT INTO source(id) VALUES (DEFAULT) RETURNING id INTO NEW.id;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER mh_user_insert
BEFORE INSERT ON mh_user
    FOR EACH ROW EXECUTE PROCEDURE mh_user_insert();

-- Delete the row from 'source' after we delete from mh_user
CREATE OR REPLACE FUNCTION mh_user_delete() RETURNS TRIGGER AS $$
    BEGIN
        DELETE FROM source WHERE id = OLD.id;
        IF NOT FOUND THEN RETURN NULL; END IF;
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER mh_user_delete
AFTER DELETE ON mh_user
    FOR EACH ROW EXECUTE PROCEDURE mh_user_delete();

-------------------------------------------------------------------------------
--
--                       Entities reliant on a "source"
--
-------------------------------------------------------------------------------

CREATE TABLE class (
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR(100) NOT NULL,
    base_class_id       INTEGER DEFAULT NULL,
    source_id           INTEGER NOT NULL,
    UNIQUE (name, source_id),
    FOREIGN KEY (base_class_id) REFERENCES class(id) ON DELETE CASCADE,
    FOREIGN KEY (source_id) REFERENCES source(id) ON DELETE CASCADE
);

CREATE TABLE spell (
    id                  SERIAL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    level               CHAR(1)     NOT NULL,
    school              VARCHAR(255) NOT NULL,
    cast_time           VARCHAR(255) NOT NULL,
    duration            VARCHAR(255) NOT NULL,
    "range"             VARCHAR(255) NOT NULL,
    comp_verbal         BOOLEAN NOT NULL,
    comp_somatic        BOOLEAN NOT NULL,
    comp_material       BOOLEAN NOT NULL,
    material_desc       TEXT,
    concentration       BOOLEAN,
    ritual              BOOLEAN,
    description         TEXT NOT NULL,
    source_id           INTEGER NOT NULL,
    FOREIGN KEY(source_id) REFERENCES source(id)
);

-- Depends only on user, we don't have cannon characters
-- mh_character because "character" is reserved, and we're much better off
-- not having to quote every time.
CREATE TABLE mh_character(
    id                     SERIAL PRIMARY KEY,
    name                   VARCHAR(255) NOT NULL,
    race                   VARCHAR(255) NOT NULL,
    spell_ability_modifier INT NULL,
    proficiency_bonus      INT NULL,
    -- alignment?
    user_id                INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES mh_user(id) ON DELETE CASCADE
);

-------------------------------------------------------------------------------
--
--                                Join tables
--
-------------------------------------------------------------------------------

CREATE TABLE character_class_levels(
    char_id             INTEGER NOT NULL,
    class_id            INTEGER NOT NULL,
    level               CHAR(2) NOT NULL, -- Max level 20
    PRIMARY KEY (char_id, class_id),
    FOREIGN KEY (char_id) REFERENCES mh_character(id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES class(id)
);

CREATE TABLE class_spells (
    spell_id            INTEGER NOT NULL,
    class_id            INTEGER NOT NULL,
    PRIMARY KEY (spell_id, class_id),
    FOREIGN KEY (spell_id) REFERENCES spell(id),
    FOREIGN KEY (class_id) REFERENCES class(id)
);

-------------------------------------------------------------------------------
-- 
--                                   Views
--
-------------------------------------------------------------------------------


CREATE VIEW cannon_classes AS
    SELECT id, name, base_class_id, source_id 
    FROM class 
    WHERE source_id IN (1, 2, 3)
;

CREATE VIEW basic_classes AS
    SELECT id, name, source_id
    FROM class
    WHERE base_class_id = NULL
;

CREATE VIEW cannon_spells AS 
    SELECT S.id, S.name, S.level, S.school, S.cast_time, S.duration, S."range", 
    S.comp_verbal, S.comp_somatic, S.comp_material, S.material_desc, 
    S.concentration, S.ritual, S.description, C.title AS source_title 
    FROM spell as S
    INNER JOIN cannon_text AS C
    ON S.source_id = C.id
;

CREATE VIEW user_spells AS
    SELECT S.id, S.name, S.level, S.school, S.cast_time, S.duration, S."range", 
    S.comp_verbal, S.comp_somatic, S.comp_material, S.material_desc, 
    S.concentration, S.ritual, S.description, U.username AS source_user 
    FROM spell as S
    INNER JOIN mh_user AS U
    ON S.source_id = U.id
;


-------------------------------------------------------------------------------
--
--                          Initial Data
--
-------------------------------------------------------------------------------

INSERT INTO cannon_text(title, abbreviation) VALUES
    ('Player''s Handbook', 'PHB'),
    ('Elemental Evil', 'EE'),
    ('Sword Coast Adventurer''s Guide', 'SCAG')
;

INSERT INTO class(id, name, base_class_id, source_id) VALUES
       (1, 'Bard', NULL, 1),
       (2, 'Cleric', NULL, 1),
       (3, 'Cleric (Arcana)', 2, 3),
       (4, 'Cleric (Knowledge)', 2, 1),
       (5, 'Cleric (Life)', 2, 1),
       (6, 'Cleric (Light)', 2, 1),
       (7, 'Cleric (Nature)', 2, 1),
       (8, 'Cleric (Tempest)', 2, 1),
       (9, 'Cleric (Trickery)', 2, 1),
       (10, 'Cleric (War)', 2, 1),
       (11, 'Cleric (Death)', 2, 1),
       (12, 'Druid', NULL, 1),
       (13, 'Druid (Arctic)', 12, 3),
       (14, 'Druid (Coast)', 12, 1),
       (15, 'Druid (Desert)', 12, 1),
       (16, 'Druid (Forest)', 12, 1),
       (17, 'Druid (Grassland)', 12, 1),
       (18, 'Druid (Mountain)', 12, 1), 
       (19, 'Druid (Swamp)', 12, 1),
       (20, 'Druid (Underdark)', 12, 1),
       (21, 'Paladin', NULL, 1),
       (22, 'Paladin (Ancients)', 21, 1),
       (23, 'Paladin (Devotion)', 21, 1),
       (24, 'Paladin (Vengeance)', 21, 1),
       (25, 'Paladin (Oathbreaker)', 21, 1),
       (26, 'Paladin (Crown)', 21, 3),
       (27, 'Ranger', NULL, 1),
       (28, 'Sorcerer', NULL, 1),
       (29, 'Warlock', NULL, 1),
       (30, 'Warlock (Archfey)', 29, 1),
       (31, 'Warlock (Fiend)', 29, 1),
       (32, 'Warlock (Great Old One)', 29, 1),
       (33, 'Warlock (Undying)', 29, 3),
       (34, 'Wizard', NULL, 1),
       (35, 'Fighter', NULL, 1),
       (36, 'Fighter (Eldritch Knight)', 35, 1),
       (37, 'Rogue', NULL, 1),
       (38, 'Rogue (Arcane Trickster)', 37, 1)
;
