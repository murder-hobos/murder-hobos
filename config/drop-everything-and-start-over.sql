SET foreign_key_checks = 0;
DROP TABLE IF EXISTS SourceText, Spell, Class, `User`, ClassSpells, CharacterLevels, `Character`, CharacterSpells;
SET foreign_key_checks = 1;

-- Create tables for murder-hobos

CREATE TABLE SourceText (
    id                  INT UNSIGNED AUTO_INCREMENT,
    `title`               VARCHAR(100) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE Spell (
    id                  INT UNSIGNED AUTO_INCREMENT,
    name                VARCHAR(50) NOT NULL,
    level               CHAR(1)     NOT NULL,
    school              VARCHAR(50) NOT NULL,
    cast_time           VARCHAR(20) NOT NULL,
    duration            VARCHAR(20) NOT NULL,
    `range`             VARCHAR(20) NOT NULL,
    comp_verbal         BOOLEAN,
    comp_somatic        BOOLEAN,
    comp_material       BOOLEAN,
    comp_material_desc  VARCHAR(255) NULL,
    concentration       BOOLEAN,
    ritual              BOOLEAN,
    source_id           INT UNSIGNED,
    PRIMARY KEY(id),
    FOREIGN KEY(source_id) REFERENCES SourceText(id)
);

CREATE TABLE User (
    id                  INT UNSIGNED AUTO_INCREMENT,
    username            VARCHAR(60) NOT NULL,
    password            CHAR(60) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE `Character`(
    id                  INT UNSIGNED AUTO_INCREMENT,
    name                VARCHAR(60) NOT NULL,
    race                VARCHAR(15) NULL,
    spell_ability_modifier INT NULL,
    proficiency_bonus      INT NULL,
    user_id                INT UNSIGNED NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES User(id) ON DELETE CASCADE
);

CREATE TABLE CharacterSpells(
    char_id             INT UNSIGNED NOT NULL,
    spell_id            INT UNSIGNED NOT NULL,
    PRIMARY KEY(char_id, spell_id),
    FOREIGN KEY(char_id) REFERENCES `Character`(id) ON DELETE CASCADE,
    FOREIGN KEY(spell_id) REFERENCES Spell(id)
);

CREATE TABLE Class (
    id                  TINYINT UNSIGNED, -- We're not going to support more than 255 classes
    name                VARCHAR(50) NOT NULL,
    base_class_id       TINYINT UNSIGNED NULL,
    PRIMARY KEY (id),
    FOREIGN KEY(base_class_id) REFERENCES Class(id)
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
