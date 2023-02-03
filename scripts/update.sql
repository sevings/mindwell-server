CREATE OR REPLACE FUNCTION to_search_string(email TEXT)
    RETURNS TEXT AS $$
BEGIN
    RETURN left(email, position('@' in email) - 1);
END
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE INDEX "index_user_search_email" ON "mindwell"."users" USING GIST
    (to_search_string("email") gist_trgm_ops);
