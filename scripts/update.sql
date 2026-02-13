-- Comment notification triggers

-- Trigger function for new comments
CREATE OR REPLACE FUNCTION notify_new_comment() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_comment', json_build_object(
        'id', NEW.id,
        'entry_id', NEW.entry_id,
        'author_id', NEW.author_id
    )::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for updated comments
CREATE OR REPLACE FUNCTION notify_update_comment() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('update_comment', json_build_object(
        'id', NEW.id,
        'entry_id', NEW.entry_id
    )::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for removed comments
CREATE OR REPLACE FUNCTION notify_remove_comment() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('remove_comment', json_build_object(
        'id', OLD.id
    )::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS comment_insert_trigger ON comments;
CREATE TRIGGER comment_insert_trigger
    AFTER INSERT ON comments
    FOR EACH ROW
    EXECUTE FUNCTION notify_new_comment();

DROP TRIGGER IF EXISTS comment_update_trigger ON comments;
CREATE TRIGGER comment_update_trigger
    AFTER UPDATE ON comments
    FOR EACH ROW
    WHEN (OLD.edit_content IS DISTINCT FROM NEW.edit_content)
    EXECUTE FUNCTION notify_update_comment();

DROP TRIGGER IF EXISTS comment_delete_trigger ON comments;
CREATE TRIGGER comment_delete_trigger
    AFTER DELETE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION notify_remove_comment();


-- Message notification triggers

-- Trigger function for new messages
CREATE OR REPLACE FUNCTION notify_new_message() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_message', json_build_object(
        'id', NEW.id,
        'chat_id', NEW.chat_id
    )::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for updated messages
CREATE OR REPLACE FUNCTION notify_update_message() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('update_message', json_build_object(
        'id', NEW.id,
        'chat_id', NEW.chat_id
    )::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger function for removed messages
CREATE OR REPLACE FUNCTION notify_remove_message() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('remove_message', json_build_object(
        'id', OLD.id,
        'chat_id', OLD.chat_id
    )::text);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS message_insert_trigger ON messages;
CREATE TRIGGER message_insert_trigger
    AFTER INSERT ON messages
    FOR EACH ROW
    EXECUTE FUNCTION notify_new_message();

DROP TRIGGER IF EXISTS message_update_trigger ON messages;
CREATE TRIGGER message_update_trigger
    AFTER UPDATE ON messages
    FOR EACH ROW
    WHEN (OLD.content IS DISTINCT FROM NEW.content OR OLD.edit_content IS DISTINCT FROM NEW.edit_content)
    EXECUTE FUNCTION notify_update_message();

DROP TRIGGER IF EXISTS message_delete_trigger ON messages;
CREATE TRIGGER message_delete_trigger
    AFTER DELETE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION notify_remove_message();


-- Message read notification trigger

-- Trigger function for message read status updates
-- When a user updates their last_read position in a chat,
-- we notify about each message that was marked as read
CREATE OR REPLACE FUNCTION notify_message_read() RETURNS TRIGGER AS $$
DECLARE
    msg_record RECORD;
    partner_name TEXT;
BEGIN
    -- Only process if last_read actually changed
    IF OLD.last_read IS DISTINCT FROM NEW.last_read THEN
        -- Find the partner's name
        SELECT users.name INTO partner_name
        FROM users
        JOIN talkers ON talkers.user_id = users.id
        WHERE talkers.chat_id = NEW.chat_id AND talkers.user_id <> NEW.user_id
        LIMIT 1;

        -- Find all messages that were just marked as read
        -- (messages sent by partner, between OLD.last_read and NEW.last_read)
        FOR msg_record IN
            SELECT id
            FROM messages
            WHERE chat_id = NEW.chat_id
              AND author_id <> NEW.user_id
              AND id > OLD.last_read
              AND id <= NEW.last_read
            ORDER BY id
        LOOP
            -- Send notification for each message
            PERFORM pg_notify('read_message', json_build_object(
                'chat_id', NEW.chat_id,
                'message_id', msg_record.id,
                'user_name', partner_name
            )::text);
        END LOOP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger on talkers table
DROP TRIGGER IF EXISTS talkers_read_trigger ON talkers;
CREATE TRIGGER talkers_read_trigger
    AFTER UPDATE ON talkers
    FOR EACH ROW
    WHEN (OLD.last_read IS DISTINCT FROM NEW.last_read)
    EXECUTE FUNCTION notify_message_read();
