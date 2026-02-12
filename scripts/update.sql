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
