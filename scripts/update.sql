CREATE OR REPLACE FUNCTION mindwell.notify_online_users() RETURNS TRIGGER AS $$
    BEGIN
        IF (NOT mindwell.is_online(OLD.last_seen_at) AND OLD.last_seen_at < NEW.last_seen_at) THEN
            PERFORM pg_notify('online_users', json_build_object(
                'id', OLD.id,
                'name', OLD.name,
                'show_name', OLD.show_name
            )::text);
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ntf_online_users
    AFTER UPDATE ON mindwell.users
    FOR EACH ROW
    EXECUTE FUNCTION mindwell.notify_online_users();
