--75_date_precision.up.sql
ALTER TABLE "scenes" ADD COLUMN "date_precision" SMALLINT;
ALTER TABLE "images" ADD COLUMN "date_precision" SMALLINT;
ALTER TABLE "galleries" ADD COLUMN "date_precision" SMALLINT;
ALTER TABLE "groups" ADD COLUMN "date_precision" SMALLINT;
ALTER TABLE "performers" ADD COLUMN "birthdate_precision" SMALLINT;
ALTER TABLE "performers" ADD COLUMN "death_date_precision" SMALLINT;

UPDATE "scenes" SET "date_precision" = 0 WHERE "date" IS NOT NULL;
UPDATE "images" SET "date_precision" = 0 WHERE "date" IS NOT NULL;
UPDATE "galleries" SET "date_precision" = 0 WHERE "date" IS NOT NULL;
UPDATE "groups" SET "date_precision" = 0 WHERE "date" IS NOT NULL;
UPDATE "performers" SET "birthdate_precision" = 0 WHERE "birthdate" IS NOT NULL;
UPDATE "performers" SET "death_date_precision" = 0 WHERE "death_date" IS NOT NULL;  
