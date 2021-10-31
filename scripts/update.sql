ALTER TABLE entries
ADD COLUMN "favorites_count" Integer DEFAULT 0 NOT NULL;

CREATE OR REPLACE FUNCTION mindwell.inc_favorites() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET favorites_count = favorites_count + 1
        WHERE id = NEW.user_id;

        UPDATE mindwell.entries
        SET favorites_count = favorites_count + 1
        WHERE id = NEW.entry_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.dec_favorites() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET favorites_count = favorites_count - 1
        WHERE id = OLD.user_id;

        UPDATE mindwell.entries
        SET favorites_count = favorites_count - 1
        WHERE id = OLD.entry_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

UPDATE entries
SET favorites_count = fav.cnt
FROM (
    SELECT entry_id, count(*) AS cnt
    FROM favorites
    GROUP BY entry_id
) fav
WHERE entries.id = fav.entry_id;
