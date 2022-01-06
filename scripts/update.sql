ALTER TABLE users
DROP COLUMN api_key;

ALTER TABLE users
DROP COLUMN valid_thru;

CREATE OR REPLACE FUNCTION mindwell.ban_user(userName TEXT) RETURNS TEXT AS $$
    BEGIN
        UPDATE mindwell.users
        SET password_hash = '', verified = false
        WHERE lower(users.name) = lower(userName);

        DELETE FROM mindwell.sessions
        WHERE user_id = (SELECT id FROM users WHERE lower(name) = lower(userName));

        RETURN (SELECT name FROM mindwell.users WHERE id = (
            SELECT invited_by FROM users WHERE lower(name) = lower(userName)
        ));
    END;
$$ LANGUAGE plpgsql;
