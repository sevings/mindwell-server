INSERT INTO "mindwell"."entry_privacy" VALUES(4, 'registered');
INSERT INTO "mindwell"."entry_privacy" VALUES(5, 'invited');
INSERT INTO "mindwell"."entry_privacy" VALUES(6, 'followers');

CREATE OR REPLACE FUNCTION mindwell.can_view_entry(user_id INTEGER, entry_id INTEGER, author_id INTEGER, entry_privacy TEXT) RETURNS BOOLEAN AS $$
    DECLARE
        allowed BOOLEAN;
    BEGIN
        IF author_id = user_id THEN
            RETURN TRUE;
        END IF;

        IF entry_privacy = 'anonymous' THEN
            RETURN user_id > 0;
        END IF;

        allowed = (SELECT can_view_tlog(user_id, author_id));

        IF NOT allowed THEN
            RETURN FALSE;
        END IF;

        CASE entry_privacy
        WHEN 'all' THEN
            RETURN TRUE;
        WHEN 'registered' THEN
            RETURN user_id > 0;
        WHEN 'invited' THEN
            IF user_id > 0 THEN
                RETURN (
                    SELECT invited_by IS NOT NULL
                    FROM users
                    WHERE users.id = user_id);
            ELSE
                RETURN FALSE;
            END IF;
        WHEN 'followers' THEN
            IF user_id > 0 THEN
                RETURN COALESCE((
                    SELECT relation.type = 'followed'
                    FROM relations
                    INNER JOIN relation ON relation.id = relations.type
                    WHERE from_id = user_id AND to_id = author_id), FALSE);
            ELSE
                RETURN FALSE;
            END IF;
        WHEN 'some' THEN
            IF user_id > 0 THEN
                SELECT TRUE
                INTO allowed
                FROM entries_privacy
                WHERE entries_privacy.user_id = can_view_entry.user_id
                    AND entries_privacy.entry_id = can_view_entry.entry_id;

                allowed = COALESCE(allowed, FALSE);
                RETURN allowed;
            ELSE
                RETURN FALSE;
            END IF;
        WHEN 'me' THEN
            RETURN FALSE;
        ELSE
            RETURN FALSE;
        END CASE;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER cnt_tlog_entries_ins ON entries;
DROP TRIGGER cnt_tlog_entries_upd_inc ON entries;
DROP TRIGGER cnt_tlog_entries_upd_dec ON entries;
DROP TRIGGER cnt_tlog_entries_del ON entries;

CREATE TRIGGER cnt_tlog_entries_ins
    AFTER INSERT ON mindwell.entries
    FOR EACH ROW
    WHEN (NEW.visible_for IN (0, 4, 5)) -- visible_for = all, registered, invited
    EXECUTE PROCEDURE mindwell.inc_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_upd_inc
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    WHEN (OLD.visible_for NOT IN (0, 4, 5) AND NEW.visible_for IN(0, 4, 5))
    EXECUTE PROCEDURE mindwell.inc_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_upd_dec
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    WHEN (OLD.visible_for IN (0, 4, 5) AND NEW.visible_for NOT IN (0, 4, 5))
    EXECUTE PROCEDURE mindwell.dec_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_del
    AFTER DELETE ON mindwell.entries
    FOR EACH ROW
    WHEN (OLD.visible_for IN (0, 4, 5))
    EXECUTE PROCEDURE mindwell.dec_tlog_entries();
