ALTER TABLE users
ADD COLUMN "pinned_entry" Integer;

ALTER TABLE users
ADD CONSTRAINT "user_pinned_entry" FOREIGN KEY("pinned_entry") REFERENCES "mindwell"."entries"("id");
