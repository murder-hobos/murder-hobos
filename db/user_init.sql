SET foreign_key_checks = 0;
DELETE FROM User;
SET foreign_key_checks = 1;
INSERT INTO `User` (id, username, password) VALUES
(1, "PHB", "totallynotsecure1"),
(2, "EE", "totallynotsecure2"),
(3, "SCAG", "totallynotsecure3")
;