ALTER TABLE mindwell.entry_images ADD COLUMN image_order INTEGER;

UPDATE mindwell.entry_images ei
SET image_order = subquery.rn
FROM (
    SELECT entry_id, image_id, ROW_NUMBER() OVER (PARTITION BY entry_id ORDER BY image_id) as rn
    FROM mindwell.entry_images
) as subquery
WHERE ei.entry_id = subquery.entry_id AND ei.image_id = subquery.image_id;

ALTER TABLE mindwell.entry_images ALTER COLUMN image_order SET NOT NULL;
