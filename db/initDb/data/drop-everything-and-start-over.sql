SET foreign_key_checks = 0;
DROP TABLE IF EXISTS `User`, Spell, Class, ClassSpells, CharacterLevels, `Character`;
SET foreign_key_checks = 1;

-- Create tables for murder-hobos

CREATE TABLE User (
    id                  INT UNSIGNED AUTO_INCREMENT,
    username            VARCHAR(60) NOT NULL,
    password            CHAR(60) NOT NULL, -- bcrypt hash length
    PRIMARY KEY(id)
);

CREATE TABLE Class (
    id                  TINYINT UNSIGNED AUTO_INCREMENT, -- We're not going to support more than 255 classes
    name                VARCHAR(50) UNIQUE NOT NULL,
    base_class_id       TINYINT UNSIGNED NULL,
    PRIMARY KEY (id),
    FOREIGN KEY(base_class_id) REFERENCES Class(id)
);

CREATE TABLE Spell (
    id                  INT UNSIGNED AUTO_INCREMENT,
    name                VARCHAR(255) NOT NULL,
    level               CHAR(1)     NOT NULL,
    school              VARCHAR(255) NOT NULL,
    cast_time           VARCHAR(255) NOT NULL,
    duration            VARCHAR(255) NOT NULL,
    `range`             VARCHAR(255) NOT NULL,
    comp_verbal         BOOLEAN NOT NULL,
    comp_somatic        BOOLEAN NOT NULL,
    comp_material       BOOLEAN NOT NULL,
    material_desc       TEXT,
    concentration       BOOLEAN,
    ritual              BOOLEAN,
    description         TEXT NOT NULL,
    source_id           INT UNSIGNED,
    PRIMARY KEY(id),
    FOREIGN KEY(source_id) REFERENCES User(id)
);

CREATE TABLE `Character`(
    id                     INT UNSIGNED AUTO_INCREMENT,
    name                   VARCHAR(255) NOT NULL,
    race                   VARCHAR(255) NULL,
    spell_ability_modifier INT NULL,
    proficiency_bonus      INT NULL,
    user_id                INT UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES User(id) ON DELETE CASCADE
);

CREATE TABLE CharacterLevels(
    char_id             INT UNSIGNED,
    class_id            TINYINT UNSIGNED,
    PRIMARY KEY (char_id, class_id),
    FOREIGN KEY (char_id) REFERENCES `Character`(id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES Class(id)
);

CREATE TABLE ClassSpells (
    spell_id            INT UNSIGNED,
    class_id            TINYINT UNSIGNED,
    PRIMARY KEY (spell_id, class_id),
    FOREIGN KEY (spell_id) REFERENCES Spell(id),
    FOREIGN KEY (class_id) REFERENCES Class(id)
);

CREATE VIEW CannonSpells AS SELECT * FROM Spell WHERE source_id IN (1, 2, 3);

-- Initialize our strong entities

INSERT INTO `User` (id, username, password) VALUES
(1, "PHB", "totallynotsecure1"),
(2, "EE", "totallynotsecure2"),
(3, "SCAG", "totallynotsecure3")
;

INSERT INTO Class(id, name, base_class_id) VALUES
       (1, "Bard", NULL),
       (2, "Cleric", NULL),
       (3, "Cleric (Arcana)", 2),
       (4, "Cleric (Knowledge)", 2),
       (5, "Cleric (Life)", 2),
       (6, "Cleric (Light)", 2),
       (7, "Cleric (Nature)", 2),
       (8, "Cleric (Tempest)", 2),
       (9, "Cleric (Trickery)", 2),
       (10, "Cleric (War)", 2),
       (11, "Cleric (Death)", 2),
       (12, "Druid", NULL),
       (13, "Druid (Arctic)", 12),
       (14, "Druid (Coast)", 12),
       (15, "Druid (Desert)", 12),
       (16, "Druid (Forest)", 12),
       (17, "Druid (Grassland)", 12),
       (18, "Druid (Mountain)", 12), 
       (19, "Druid (Swamp)", 12),
       (20, "Druid (Underdark)", 12),
       (21, "Paladin", NULL),
       (22, "Paladin (Ancients)", 21),
       (23, "Paladin (Devotion)", 21),
       (24, "Paladin (Vengeance)", 21),
       (25, "Paladin (Oathbreaker)", 21),
       (26, "Paladin (Crown)", 21),
       (27, "Ranger", NULL),
       (28, "Sorcerer", NULL),
       (29, "Warlock", NULL),
       (30, "Warlock (Archfey)", 29),
       (31, "Warlock (Fiend)", 29),
       (32, "Warlock (Great Old One)", 29),
       (33, "Warlock (Undying)", 29),
       (34, "Wizard", NULL),
       (35, "Fighter", NULL),
       (36, "Fighter (Eldritch Knight)", 35),
       (37, "Rogue", NULL),
       (38, "Rogue (Arcane Trickster)", 37)
;