INSERT INTO "mindwell"."complain_type" VALUES(3, 'user');
INSERT INTO "mindwell"."complain_type" VALUES(4, 'theme');

ALTER TABLE "mindwell"."adm"
ADD COLUMN "tracking" Text NOT NULL DEFAULT '';

ALTER TABLE "mindwell"."adm"
ADD COLUMN "grandfather_comment" Text NOT NULL DEFAULT '';

ALTER TABLE "mindwell"."entries"
ADD COLUMN "is_shared" Boolean DEFAULT FALSE NOT NULL;
