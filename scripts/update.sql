CREATE OR REPLACE FUNCTION mindwell.notify_entries() RETURNS TRIGGER AS $$
    BEGIN
        IF OLD.author_id <> NEW.author_id THEN
            PERFORM pg_notify('moved_entries', json_build_object(
                'id', OLD.id,
                'title', OLD.title,
                'author', json_build_object(
                    'id', OLD.author_id
                ),
                'user', json_build_object(
                    'id', OLD.user_id
                )
            )::text);
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ntf_entries
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    EXECUTE FUNCTION mindwell.notify_entries();

ALTER TABLE users
ADD COLUMN "email_moved_entries" Boolean NOT NULL DEFAULT FALSE;

UPDATE users
SET email_moved_entries = TRUE
WHERE email_invites = TRUE;

ALTER TABLE users
ADD COLUMN "telegram_moved_entries" Boolean NOT NULL DEFAULT TRUE;

UPDATE users
SET telegram_moved_entries = FALSE
WHERE telegram_invites = FALSE;

INSERT INTO "mindwell"."notification_type" VALUES(12, 'entry_moved');
