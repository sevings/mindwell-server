ALTER TABLE "mindwell"."adm"
ADD COLUMN "phone" Text NOT NULL DEFAULT '';

ALTER TABLE "mindwell"."adm"
ALTER COLUMN "phone" DROP DEFAULT;
