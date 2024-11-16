ALTER TABLE users
ADD COLUMN "alt_of" Text;

ALTER TABLE users
ADD COLUMN "confirmed_alt" Boolean DEFAULT FALSE NOT NULL;
