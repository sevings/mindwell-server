UPDATE entries
SET visible_for = (SELECT id FROM entry_privacy WHERE type = 'followers')
WHERE visible_for = (SELECT id FROM entry_privacy WHERE type = 'some');

CREATE TABLE "mindwell"."authority" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."authority" VALUES(0, 'user');
INSERT INTO "mindwell"."authority" VALUES(1, 'admin');

ALTER TABLE users
ADD COLUMN "authority" Integer DEFAULT 0 NOT NULL;

ALTER TABLE users
ADD CONSTRAINT "enum_user_authority" FOREIGN KEY("authority") REFERENCES "mindwell"."authority"("id");

UPDATE users
SET authority = 1
WHERE lower(name) = 'mindwell';
