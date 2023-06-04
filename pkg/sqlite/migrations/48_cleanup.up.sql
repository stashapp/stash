-- Cleanup old invalid dates
UPDATE `scenes` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `galleries` SET `date` = NULL WHERE `date` = '0001-01-01' OR `date` = '';
UPDATE `performers` SET `birthdate` = NULL WHERE `birthdate` = '0001-01-01' OR `birthdate` = '';
UPDATE `performers` SET `death_date` = NULL WHERE `death_date` = '0001-01-01' OR `death_date` = '';
