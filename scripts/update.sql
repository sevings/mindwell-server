ALTER TABLE mindwell.images
ADD COLUMN "preview_extension" Text DEFAULT '' NOT NULL;

UPDATE mindwell.images
SET preview_extension = 'jpg'
WHERE extension = 'gif';
