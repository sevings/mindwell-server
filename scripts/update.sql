-- INSERT INTO "mindwell"."authority" VALUES(2, 'moderator');

ALTER TABLE "mindwell"."users"
ADD COLUMN "complain_ban" Date DEFAULT CURRENT_DATE NOT NULL;
