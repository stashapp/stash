UPDATE `scenes` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `galleries` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `images` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `movies` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `performers` SET `rating` = (`rating` * 20) WHERE `rating` < 6;
UPDATE `studios` SET `rating` = (`rating` * 20) WHERE `rating` < 6;