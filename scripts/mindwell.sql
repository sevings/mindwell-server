SET client_encoding = 'UTF8';

CREATE SCHEMA "mindwell";

CREATE EXTENSION IF NOT EXISTS pg_trgm; -- search users
CREATE EXTENSION IF NOT EXISTS rum;     -- search entries

ALTER DATABASE mindwell SET search_path TO mindwell, public;


CREATE TABLE "mindwell"."gender" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."gender" VALUES(0, 'not set');
INSERT INTO "mindwell"."gender" VALUES(1, 'male');
INSERT INTO "mindwell"."gender" VALUES(2, 'female');


CREATE TABLE "mindwell"."user_privacy" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."user_privacy" VALUES(0, 'all');
INSERT INTO "mindwell"."user_privacy" VALUES(1, 'followers');
INSERT INTO "mindwell"."user_privacy" VALUES(2, 'invited');
INSERT INTO "mindwell"."user_privacy" VALUES(3, 'registered');


CREATE TABLE "mindwell"."user_chat_privacy" (
                                           "id" Integer UNIQUE NOT NULL,
                                           "type" Text NOT NULL );

INSERT INTO "mindwell"."user_chat_privacy" VALUES(0, 'invited');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(1, 'followers');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(2, 'friends');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(3, 'me');


CREATE TABLE "mindwell"."authority" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."authority" VALUES(0, 'user');
INSERT INTO "mindwell"."authority" VALUES(1, 'admin');
INSERT INTO "mindwell"."authority" VALUES(2, 'moderator');


CREATE TABLE "mindwell"."font_family" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."font_family" VALUES(0, 'Arial');


CREATE TABLE "mindwell"."alignment" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."alignment" VALUES(0, 'left');
INSERT INTO "mindwell"."alignment" VALUES(1, 'right');
INSERT INTO "mindwell"."alignment" VALUES(2, 'center');
INSERT INTO "mindwell"."alignment" VALUES(3, 'justify');


CREATE OR REPLACE FUNCTION to_search_string(name TEXT, show_name TEXT, country TEXT, city TEXT)
   RETURNS TEXT AS $$
BEGIN
  RETURN name || ' ' || show_name || ' ' || country || ' ' || city;
END
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION to_search_string(email TEXT)
    RETURNS TEXT AS $$
BEGIN
    RETURN left(email, position('@' in email) - 1);
END
$$ LANGUAGE plpgsql IMMUTABLE;


CREATE TABLE "mindwell"."users" (
	"id" Serial NOT NULL,
	"name" Text NOT NULL,
	"show_name" Text DEFAULT '' NOT NULL,
	"password_hash" Bytea NOT NULL,
	"gender" Integer DEFAULT 0 NOT NULL,
	"is_daylog" Boolean DEFAULT false NOT NULL,
    "show_in_tops" Boolean DEFAULT true NOT NULL,
    "privacy" Integer DEFAULT 0 NOT NULL,
    "chat_privacy" Integer DEFAULT 0 NOT NULL,
	"title" Text DEFAULT '' NOT NULL,
	"last_seen_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "rank" Integer NOT NULL,
	"karma" Real DEFAULT 0 NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"last_invite" Date DEFAULT CURRENT_DATE NOT NULL,
    "invited_by" Integer,
	"birthday" Date,
	"css" Text DEFAULT '' NOT NULL,
	"entries_count" Integer DEFAULT 0 NOT NULL,
	"followings_count" Integer DEFAULT 0 NOT NULL,
	"followers_count" Integer DEFAULT 0 NOT NULL,
	"comments_count" Integer DEFAULT 0 NOT NULL,
	"ignored_count" Integer DEFAULT 0 NOT NULL,
	"invited_count" Integer DEFAULT 0 NOT NULL,
	"favorites_count" Integer DEFAULT 0 NOT NULL,
	"tags_count" Integer DEFAULT 0 NOT NULL,
	"badges_count" Integer DEFAULT 0 NOT NULL,
	"country" Text DEFAULT '' NOT NULL,
	"city" Text DEFAULT '' NOT NULL,
	"email" Text NOT NULL,
	"verified" Boolean DEFAULT false NOT NULL,
	"avatar" Text DEFAULT '' NOT NULL,
	"cover" Text DEFAULT '' NOT NULL,
	"font_family" Integer DEFAULT 0 NOT NULL,
	"font_size" SmallInt DEFAULT 100 NOT NULL,
	"text_alignment" Integer DEFAULT 0 NOT NULL,
	"text_color" Character( 7 ) DEFAULT '#000000' NOT NULL,
	"background_color" Character( 7 ) DEFAULT '#ffffff' NOT NULL,
    "email_comments" Boolean NOT NULL DEFAULT FALSE,
    "email_followers" Boolean NOT NULL DEFAULT FALSE,
    "email_invites" Boolean NOT NULL DEFAULT FALSE,
    "email_moved_entries" Boolean NOT NULL DEFAULT FALSE,
    "email_badges" Boolean NOT NULL DEFAULT FALSE,
    "telegram_comments" Boolean NOT NULL DEFAULT TRUE,
    "telegram_followers" Boolean NOT NULL DEFAULT TRUE,
    "telegram_invites" Boolean NOT NULL DEFAULT TRUE,
    "telegram_moved_entries" Boolean NOT NULL DEFAULT TRUE,
    "telegram_badges" Boolean NOT NULL DEFAULT TRUE,
    "telegram_messages" Boolean NOT NULL DEFAULT TRUE,
    "send_wishes" Boolean NOT NULL DEFAULT TRUE,
    "invite_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "vote_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "comment_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "live_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "complain_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "user_ban" Date DEFAULT CURRENT_DATE NOT NULL,
    "adm_ban" Boolean DEFAULT TRUE NOT NULL,
    "shadow_ban" Boolean DEFAULT FALSE NOT NULL,
    "telegram" Integer,
    "authority" Integer DEFAULT 0 NOT NULL,
    "creator_id" Integer,
    "pinned_entry" Integer,
    "alt_of" Text,
    "confirmed_alt" Boolean DEFAULT FALSE NOT NULL,
	CONSTRAINT "unique_user_id" PRIMARY KEY( "id" ),
    CONSTRAINT "enum_user_gender" FOREIGN KEY("gender") REFERENCES "mindwell"."gender"("id"),
    CONSTRAINT "enum_user_privacy" FOREIGN KEY("privacy") REFERENCES "mindwell"."user_privacy"("id"),
    CONSTRAINT "enum_user_chat_privacy" FOREIGN KEY("chat_privacy") REFERENCES "mindwell"."user_chat_privacy"("id"),
    CONSTRAINT "enum_user_alignment" FOREIGN KEY("text_alignment") REFERENCES "mindwell"."alignment"("id"),
    CONSTRAINT "enum_user_font_family" FOREIGN KEY("font_family") REFERENCES "mindwell"."font_family"("id"),
    CONSTRAINT "enum_user_authority" FOREIGN KEY("authority") REFERENCES "mindwell"."authority"("id"),
    CONSTRAINT "theme_creator" FOREIGN KEY("creator_id") REFERENCES "mindwell"."users"("id") );
 ;


CREATE UNIQUE INDEX "index_user_id" ON "mindwell"."users" USING btree( "id" );

CREATE UNIQUE INDEX "index_user_name" ON "mindwell"."users" USING btree( lower("name") );

CREATE UNIQUE INDEX "index_user_email" ON "mindwell"."users" USING btree( lower("email") );

CREATE UNIQUE INDEX "index_telegram" ON "mindwell"."users" USING btree( "telegram" );

CREATE INDEX "index_invited_by" ON "mindwell"."users" USING btree( "invited_by" );

CREATE INDEX "index_user_rank" ON "mindwell"."users" USING btree( "rank" );

CREATE INDEX "index_user_search" ON "mindwell"."users" USING GIST
    (to_search_string("name", "show_name", "country", "city") gist_trgm_ops);

CREATE INDEX "index_user_search_email" ON "mindwell"."users" USING GIST
    (to_search_string("email") gist_trgm_ops);

CREATE OR REPLACE FUNCTION mindwell.count_invited_upd() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET invited_count = invited_count + 1
        WHERE id = NEW.invited_by;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_invited_upd
    AFTER UPDATE ON mindwell.users
    FOR EACH ROW
    WHEN (OLD.invited_by IS NULL AND NEW.invited_by IS NOT NULL)
    EXECUTE PROCEDURE mindwell.count_invited_upd();

CREATE OR REPLACE FUNCTION mindwell.count_invited_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET invited_count = invited_count - 1
        WHERE id = OLD.invited_by;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_invited_del
    AFTER DELETE ON mindwell.users
    FOR EACH ROW
    WHEN (OLD.invited_by IS NOT NULL)
    EXECUTE PROCEDURE mindwell.count_invited_del();

CREATE OR REPLACE FUNCTION mindwell.allow_adm_upd() RETURNS TRIGGER AS $$
    BEGIN
        NEW.adm_ban = false;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER alw_adm_upd
    BEFORE UPDATE ON mindwell.users
    FOR EACH ROW
    WHEN (OLD.invited_by IS NULL AND NEW.invited_by IS NOT NULL)
    EXECUTE PROCEDURE mindwell.allow_adm_upd();

CREATE  OR REPLACE FUNCTION mindwell.is_online(last_seen_at TIMESTAMP WITH TIME ZONE) RETURNS BOOLEAN AS $$
    BEGIN
        RETURN now() - last_seen_at < interval '5 minutes';
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.user_age(birthday date) RETURNS integer AS $$
    BEGIN
        RETURN extract(year from age(birthday))::integer;
    END;
$$ LANGUAGE plpgsql;

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


CREATE TABLE "mindwell"."adm" (
	"name" Text NOT NULL,
    "fullname" Text NOT NULL,
    "postcode" Text NOT NULL,
    "country" Text NOT NULL,
    "address" Text NOT NULL,
    "phone" Text NOT NULL,
    "comment" Text NOT NULL,
    "anonymous" Boolean NOT NULL,
    "grandfather" Text NOT NULL DEFAULT '',
    "sent" Boolean NOT NULL DEFAULT false,
    "received" Boolean NOT NULL DEFAULT false,
    "tracking" Text NOT NULL DEFAULT '',
    "grandfather_comment" Text NOT NULL DEFAULT '');
;

CREATE UNIQUE INDEX "index_adm" ON "mindwell"."adm" USING btree( lower("name") );

CREATE OR REPLACE FUNCTION mindwell.ban_adm() RETURNS VOID AS $$
    UPDATE mindwell.users
    SET adm_ban = true
    WHERE name IN (
        SELECT gs.name
        FROM mindwell.adm AS gs
        JOIN mindwell.adm AS gf ON gf.grandfather = gs.name
        WHERE (NOT gf.sent AND NOT gf.received) OR (gs.sent AND NOT gs.received)
    );
$$ LANGUAGE SQL;


CREATE TABLE "mindwell"."wish_states" (
    "id" Integer UNIQUE NOT NULL,
    "state" Text NOT NULL
);

INSERT INTO "mindwell"."wish_states"(id, state) VALUES (0, 'new');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (1, 'sent');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (2, 'declined');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (3, 'complained');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (4, 'thanked');


CREATE TABLE "mindwell"."wishes" (
    "id" Serial NOT NULL,
    "from_id" Integer NOT NULL,
    "to_id" Integer NOT NULL,
    "content" Text DEFAULT '' NOT NULL,
    "state" Integer DEFAULT 0 NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "unique_wish_id" PRIMARY KEY( "id" ),
    CONSTRAINT "wish_sender" FOREIGN KEY("from_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "wish_receiver" FOREIGN KEY("to_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_wish_state" FOREIGN KEY("state") REFERENCES "mindwell"."wish_states"("id")
);

CREATE INDEX "index_wish_id" ON "mindwell"."wishes" USING btree( "id" );

CREATE INDEX "index_wish_from_id" ON "mindwell"."wishes" USING btree( "from_id" );

CREATE INDEX "index_wish_to_id" ON "mindwell"."wishes" USING btree( "to_id" );


CREATE TABLE "mindwell"."invite_words" (
    "id" Serial NOT NULL,
    "word" Text NOT NULL,
	CONSTRAINT "unique_word_id" PRIMARY KEY( "id" ),
	CONSTRAINT "unique_word" UNIQUE( "word" ) );
;

CREATE UNIQUE INDEX "index_invite_word_id" ON "mindwell"."invite_words" USING btree( "id" );

CREATE UNIQUE INDEX "index_invite_word" ON "mindwell"."invite_words" USING btree( "word" );

INSERT INTO mindwell.invite_words ("word") VALUES('acknown');
INSERT INTO mindwell.invite_words ("word") VALUES('aery');
INSERT INTO mindwell.invite_words ("word") VALUES('affectioned');
INSERT INTO mindwell.invite_words ("word") VALUES('agnize');
INSERT INTO mindwell.invite_words ("word") VALUES('ambition');


CREATE TABLE "mindwell"."invites" (
    "id" Serial NOT NULL,
    "referrer_id" Integer NOT NULL,
    "word1" Integer NOT NULL,
    "word2" Integer NOT NULL,
    "word3" Integer NOT NULL,
	"created_at" Date DEFAULT CURRENT_DATE NOT NULL,
    CONSTRAINT "unique_invite_id" PRIMARY KEY( "id" ),
    CONSTRAINT "invite_word1" FOREIGN KEY("word1") REFERENCES "mindwell"."invite_words"("id"),
    CONSTRAINT "invite_word2" FOREIGN KEY("word2") REFERENCES "mindwell"."invite_words"("id"),
    CONSTRAINT "invite_word3" FOREIGN KEY("word3") REFERENCES "mindwell"."invite_words"("id") );
;

CREATE INDEX "index_referrer_id" ON "mindwell"."invites" USING btree( "referrer_id" );

CREATE UNIQUE INDEX "index_invite_words" ON "mindwell"."invites" USING btree( "word1", "word2", "word3" );

INSERT INTO mindwell.invites (referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1);
INSERT INTO mindwell.invites (referrer_id, word1, word2, word3) VALUES(1, 2, 2, 2);
INSERT INTO mindwell.invites (referrer_id, word1, word2, word3) VALUES(1, 3, 3, 3);

CREATE OR REPLACE FUNCTION give_invite(userName TEXT) RETURNS VOID AS $$
    DECLARE
        wordCount INTEGER;
        userId INTEGER;
    BEGIN
        wordCount = (SELECT COUNT(*) FROM invite_words);
        userId = (SELECT id FROM users WHERE lower(name) = lower(userName));

        INSERT INTO invites(referrer_id, word1, word2, word3)
            VALUES(userId,
                ceil(random() * wordCount),
                ceil(random() * wordCount),
                ceil(random() * wordCount));
    END;
$$ LANGUAGE plpgsql;

CREATE VIEW mindwell.unwrapped_invites AS
SELECT invites.id AS id,
    users.id AS user_id,
    lower(users.name) AS name,
    one.word AS word1,
    two.word AS word2,
    three.word AS word3
FROM mindwell.invites, mindwell.users,
    mindwell.invite_words AS one,
    mindwell.invite_words AS two,
    mindwell.invite_words AS three
WHERE invites.referrer_id = users.id
    AND invites.word1 = one.id
    AND invites.word2 = two.id
    AND invites.word3 = three.id;


CREATE TABLE "mindwell"."relation" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."relation" VALUES(0, 'followed');
INSERT INTO "mindwell"."relation" VALUES(1, 'requested');
INSERT INTO "mindwell"."relation" VALUES(2, 'cancelled');
INSERT INTO "mindwell"."relation" VALUES(3, 'ignored');
INSERT INTO "mindwell"."relation" VALUES(4, 'hidden');


CREATE TABLE "mindwell"."relations" (
	"from_id" Integer NOT NULL,
	"to_id" Integer NOT NULL,
	"type" Integer NOT NULL,
	"changed_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT "unique_relation" PRIMARY KEY ("from_id" , "to_id"),
    CONSTRAINT "unique_from_relation" FOREIGN KEY ("from_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "unique_to_relation" FOREIGN KEY ("to_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_relation_type" FOREIGN KEY("type") REFERENCES "mindwell"."relation"("id") );
 ;

CREATE INDEX "index_related_to_users" ON "mindwell"."relations" USING btree( "to_id" );

CREATE INDEX "index_related_from_users" ON "mindwell"."relations" USING btree( "from_id" );

CREATE OR REPLACE FUNCTION mindwell.count_relations_ins() RETURNS TRIGGER AS $$
    BEGIN
        IF (NEW."type" = (SELECT id FROM mindwell.relation WHERE "type" = 'followed')) THEN
            UPDATE mindwell.users
            SET followers_count = followers_count + 1
            WHERE id = NEW.to_id;
            UPDATE mindwell.users
            SET followings_count = followings_count + 1
            WHERE id = NEW.from_id;
        ELSIF (NEW."type" = (SELECT id FROM mindwell.relation WHERE "type" = 'ignored')) THEN
            UPDATE mindwell.users
            SET ignored_count = ignored_count + 1
            WHERE id = NEW.from_id;
        END IF;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_relations_ins
    AFTER INSERT OR UPDATE ON mindwell.relations
    FOR EACH ROW EXECUTE PROCEDURE mindwell.count_relations_ins();

CREATE OR REPLACE FUNCTION mindwell.count_relations_del() RETURNS TRIGGER AS $$
    BEGIN
        IF (OLD."type" = (SELECT id FROM mindwell.relation WHERE "type" = 'followed')) THEN
            UPDATE mindwell.users
            SET followers_count = followers_count - 1
            WHERE id = OLD.to_id;
            UPDATE mindwell.users
            SET followings_count = followings_count - 1
            WHERE id = OLD.from_id;
        ELSIF (OLD."type" = (SELECT id FROM mindwell.relation WHERE "type" = 'ignored')) THEN
            UPDATE users
            SET ignored_count = ignored_count - 1
            WHERE id = OLD.from_id;
        END IF;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_relations_del
    AFTER UPDATE OR DELETE ON mindwell.relations
    FOR EACH ROW EXECUTE PROCEDURE mindwell.count_relations_del();

CREATE OR REPLACE FUNCTION mindwell.del_relation_from_ignored() RETURNS TRIGGER AS $$
    DECLARE
        ignored Integer;
    BEGIN
        ignored = (SELECT id FROM mindwell.relation WHERE "type" = 'ignored');

        IF (NEW."type" = ignored) THEN
            DELETE FROM relations
            WHERE relations.from_id = NEW.to_id
                AND relations.to_id = NEW.from_id
                AND relations."type" != ignored;
        END IF;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER relation_from_ignored
    AFTER INSERT OR UPDATE ON mindwell.relations
    FOR EACH ROW EXECUTE PROCEDURE mindwell.del_relation_from_ignored();


CREATE TABLE "mindwell"."entry_privacy" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."entry_privacy" VALUES(0, 'all');
INSERT INTO "mindwell"."entry_privacy" VALUES(1, 'some');
INSERT INTO "mindwell"."entry_privacy" VALUES(2, 'me');
INSERT INTO "mindwell"."entry_privacy" VALUES(3, 'anonymous');
INSERT INTO "mindwell"."entry_privacy" VALUES(4, 'registered');
INSERT INTO "mindwell"."entry_privacy" VALUES(5, 'invited');
INSERT INTO "mindwell"."entry_privacy" VALUES(6, 'followers');


CREATE TABLE "mindwell"."categories" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."categories" VALUES(0, 'tweet');
INSERT INTO "mindwell"."categories" VALUES(1, 'longread');
INSERT INTO "mindwell"."categories" VALUES(2, 'media');
INSERT INTO "mindwell"."categories" VALUES(3, 'comment');


CREATE OR REPLACE FUNCTION to_search_vector(title TEXT, content TEXT)
   RETURNS tsvector AS $$
BEGIN
  RETURN to_tsvector(title || '\n' || content);
END
$$ LANGUAGE plpgsql IMMUTABLE;


CREATE TABLE "mindwell"."entries" (
	"id" Serial NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"author_id" Integer NOT NULL,
	"user_id" Integer NOT NULL,
	"title" Text DEFAULT '' NOT NULL,
    "edit_content" Text NOT NULL,
	"word_count" Integer NOT NULL,
	"visible_for" Integer NOT NULL,
	"is_votable" Boolean NOT NULL,
	"is_commentable" Boolean DEFAULT TRUE NOT NULL,
    "is_anonymous" Boolean DEFAULT FALSE NOT NULL,
    "is_shared" Boolean DEFAULT FALSE NOT NULL,
    "in_live" Boolean DEFAULT TRUE NOT NULL,
	"rating" Real DEFAULT 0 NOT NULL,
    "up_votes" Integer DEFAULT 0 NOT NULL,
    "down_votes" Integer DEFAULT 0 NOT NULL,
    "vote_sum" Real DEFAULT 0 NOT NULL,
    "weight_sum" Real DEFAULT 0 NOT NULL,
    "category" Integer NOT NULL,
	"comments_count" Integer DEFAULT 0 NOT NULL,
	"favorites_count" Integer DEFAULT 0 NOT NULL,
    "last_comment" Integer,
	CONSTRAINT "unique_entry_id" PRIMARY KEY( "id" ),
    CONSTRAINT "entry_author_id" FOREIGN KEY("author_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "entry_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_entry_privacy" FOREIGN KEY("visible_for") REFERENCES "mindwell"."entry_privacy"("id"),
    CONSTRAINT "entry_category" FOREIGN KEY("category") REFERENCES "mindwell"."categories"("id") );
 ;

CREATE INDEX "index_entry_id" ON "mindwell"."entries" USING btree( "id" );

CREATE INDEX "index_entry_date" ON "mindwell"."entries" USING btree( "created_at" );

CREATE INDEX "index_entry_author_id" ON "mindwell"."entries" USING btree( "author_id" );

CREATE INDEX "index_entry_user_id" ON "mindwell"."entries" USING btree( "user_id" );

CREATE INDEX "index_entry_rating" ON "mindwell"."entries" USING btree( "rating" );

CREATE INDEX "index_entry_word_count" ON "mindwell"."entries" USING btree( "word_count" );

CREATE INDEX "index_last_comment_id" ON "mindwell"."entries" USING btree( "last_comment" );

CREATE INDEX "index_entry_search" ON "mindwell"."entries" USING rum
    (to_search_vector("title", "edit_content") rum_tsvector_ops);

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

CREATE OR REPLACE FUNCTION mindwell.inc_tlog_entries() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET entries_count = entries_count + 1
        WHERE id = NEW.author_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.dec_tlog_entries() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET entries_count = entries_count - 1
        WHERE id = OLD.author_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ntf_entries
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    EXECUTE FUNCTION mindwell.notify_entries();

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


ALTER TABLE "mindwell"."users"
ADD CONSTRAINT "user_pinned_entry" FOREIGN KEY("pinned_entry") REFERENCES "mindwell"."entries"("id");


CREATE TABLE "mindwell"."tags" (
    "id" Serial NOT NULL,
    "tag" Text NOT NULL,
    CONSTRAINT "unique_tag_id" PRIMARY KEY( "id" ) );
;

CREATE UNIQUE INDEX "index_tag" ON "mindwell"."tags" USING btree( "tag" ) ;

CREATE INDEX "index_tag_search" ON "mindwell"."tags" USING GIST("tag" gist_trgm_ops);


CREATE TABLE "mindwell"."entry_tags" (
    "entry_id" Integer NOT NULL,
    "tag_id" Integer NOT NULL,
    CONSTRAINT "entry_tags_entry" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "entry_tags_tag" FOREIGN KEY("tag_id") REFERENCES "mindwell"."tags"("id"),
    CONSTRAINT "unique_entry_tag" UNIQUE("entry_id", "tag_id") );
;

CREATE INDEX "index_entry_tags_entry" ON "mindwell"."entry_tags" USING btree( "entry_id" );

CREATE INDEX "index_entry_tags_tag" ON "mindwell"."entry_tags" USING btree( "tag_id" );

CREATE OR REPLACE FUNCTION mindwell.count_tags_ins() RETURNS TRIGGER AS $$
    BEGIN
        WITH authors AS
        (
            SELECT DISTINCT author_id as id
            FROM changes
            INNER JOIN mindwell.entries ON changes.entry_id = entries.id
        )
        UPDATE mindwell.users
        SET tags_count = counts.cnt
        FROM
        (
            SELECT authors.id, COUNT(DISTINCT tag_id) as cnt
            FROM authors
            LEFT JOIN mindwell.entries ON authors.id = entries.author_id
            LEFT JOIN mindwell.entry_privacy ON entries.visible_for = entry_privacy.id
            LEFT JOIN mindwell.entry_tags ON entries.id = entry_tags.entry_id
            WHERE (entry_privacy.type = 'all' OR entry_privacy.type IS NULL)
            GROUP BY authors.id
        ) AS counts
        WHERE counts.id = users.id
            AND tags_count <> counts.cnt;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_tags_ins
    AFTER INSERT ON mindwell.entry_tags
    REFERENCING NEW TABLE as changes
    FOR EACH STATEMENT EXECUTE PROCEDURE mindwell.count_tags_ins();

CREATE OR REPLACE FUNCTION mindwell.count_tags_upd() RETURNS TRIGGER AS $$
    BEGIN
        WITH authors AS
        (
            SELECT OLD.author_id AS id
        )
        UPDATE mindwell.users
        SET tags_count = counts.cnt
        FROM
        (
            SELECT authors.id, COUNT(DISTINCT tag_id) as cnt
            FROM authors
            LEFT JOIN mindwell.entries ON authors.id = entries.author_id
            LEFT JOIN mindwell.entry_privacy ON entries.visible_for = entry_privacy.id
            LEFT JOIN mindwell.entry_tags ON entries.id = entry_tags.entry_id
            WHERE (entry_privacy.type = 'all' OR entry_privacy.type IS NULL)
            GROUP BY authors.id
        ) AS counts
        WHERE users.id = counts.id
            AND tags_count <> counts.cnt;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_tags_upd
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    WHEN ( OLD.visible_for <> NEW.visible_for )
    EXECUTE PROCEDURE mindwell.count_tags_upd();

CREATE OR REPLACE FUNCTION mindwell.count_tags_del() RETURNS TRIGGER AS $$
    BEGIN
        WITH authors AS
        (
            SELECT DISTINCT author_id as id
            FROM changes
        )
        UPDATE mindwell.users
        SET tags_count = counts.cnt
        FROM
        (
            SELECT authors.id, COUNT(DISTINCT tag_id) as cnt
            FROM authors
            LEFT JOIN mindwell.entries ON authors.id = entries.author_id
            LEFT JOIN mindwell.entry_privacy ON entries.visible_for = entry_privacy.id
            LEFT JOIN mindwell.entry_tags ON entries.id = entry_tags.entry_id
            WHERE (entry_privacy.type = 'all' OR entry_privacy.type IS NULL)
            GROUP BY authors.id
        ) AS counts
        WHERE users.id = counts.id
            AND tags_count <> counts.cnt;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_tags_del
    AFTER DELETE ON mindwell.entries
    REFERENCING OLD TABLE AS changes
    FOR EACH STATEMENT
    EXECUTE PROCEDURE mindwell.count_tags_del();


CREATE TABLE "mindwell"."favorites" (
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    "date" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "favorite_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "favorite_entry_id" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_user_favorite" UNIQUE("user_id", "entry_id") );
 ;

CREATE INDEX "index_favorite_entries" ON "mindwell"."favorites" USING btree( "entry_id" );

CREATE INDEX "index_favorite_users" ON "mindwell"."favorites" USING btree( "user_id" );

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

CREATE TRIGGER cnt_favorites_inc
    AFTER INSERT ON mindwell.favorites
    FOR EACH ROW EXECUTE PROCEDURE mindwell.inc_favorites();

CREATE TRIGGER cnt_favorites_dec
    AFTER DELETE ON mindwell.favorites
    FOR EACH ROW EXECUTE PROCEDURE mindwell.dec_favorites();


CREATE TABLE "mindwell"."watching" (
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    "date" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "watching_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "watching_entry_id" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_user_watching" UNIQUE("user_id", "entry_id") );
 ;

CREATE INDEX "index_watching_entries" ON "mindwell"."watching" USING btree( "entry_id" );

CREATE INDEX "index_watching_users" ON "mindwell"."watching" USING btree( "user_id" );


CREATE TABLE "mindwell"."entry_votes" (
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    "vote" Real NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "entry_vote_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "entry_vote_entry_id" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_entry_vote" UNIQUE("user_id", "entry_id") );
 ;

CREATE INDEX "index_voted_entries" ON "mindwell"."entry_votes" USING btree( "entry_id" );

CREATE INDEX "index_voted_users" ON "mindwell"."entry_votes" USING btree( "user_id" );

CREATE OR REPLACE FUNCTION mindwell.entry_votes_ins() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.entries
        SET up_votes = up_votes + (NEW.vote > 0)::int,
            down_votes = down_votes + (NEW.vote < 0)::int,
            vote_sum = vote_sum + NEW.vote,
            weight_sum = weight_sum + abs(NEW.vote),
            rating = atan2(weight_sum + abs(NEW.vote), 2)
                * (vote_sum + NEW.vote) / (weight_sum + abs(NEW.vote)) / pi() * 200
        WHERE id = NEW.entry_id;

        WITH entry AS (
            SELECT user_id, category
            FROM mindwell.entries
            WHERE id = NEW.entry_id
        )
        UPDATE mindwell.vote_weights
        SET vote_count = vote_count + 1,
            vote_sum = vote_sum + NEW.vote,
            weight_sum = weight_sum + abs(NEW.vote),
            weight = atan2(vote_count + 1, 200) * (vote_sum + NEW.vote)
                / (weight_sum + abs(NEW.vote)) / pi() * 2
        FROM entry
        WHERE vote_weights.user_id = entry.user_id
            AND vote_weights.category = entry.category;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.entry_votes_upd() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.entries
        SET up_votes = up_votes - (OLD.vote > 0)::int + (NEW.vote > 0)::int,
            down_votes = down_votes - (OLD.vote < 0)::int + (NEW.vote < 0)::int,
            vote_sum = vote_sum - OLD.vote + NEW.vote,
            weight_sum = weight_sum - abs(OLD.vote) + abs(NEW.vote),
            rating = atan2(weight_sum - abs(OLD.vote) + abs(NEW.vote), 2)
                * (vote_sum - OLD.vote + NEW.vote) / (weight_sum - abs(OLD.vote) + abs(NEW.vote)) / pi() * 200
        WHERE id = NEW.entry_id;

        WITH entry AS (
            SELECT user_id, category
            FROM mindwell.entries
            WHERE id = NEW.entry_id
        )
        UPDATE mindwell.vote_weights
        SET vote_sum = vote_sum - OLD.vote + NEW.vote,
            weight_sum = weight_sum - abs(OLD.vote) + abs(NEW.vote),
            weight = atan2(vote_count, 200) * (vote_sum - OLD.vote + NEW.vote)
                / (weight_sum - abs(OLD.vote) + abs(NEW.vote)) / pi() * 2
        FROM entry
        WHERE vote_weights.user_id = entry.user_id
            AND vote_weights.category = entry.category;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.entry_votes_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.entries
        SET up_votes = up_votes - (OLD.vote > 0)::int,
            down_votes = down_votes - (OLD.vote < 0)::int,
            vote_sum = vote_sum - OLD.vote,
            weight_sum = weight_sum - abs(OLD.vote),
            rating = CASE WHEN weight_sum = abs(OLD.vote) THEN 0
                ELSE atan2(weight_sum - abs(OLD.vote), 2)
                    * (vote_sum - OLD.vote) / (weight_sum - abs(OLD.vote)) / pi() * 200
                END
        WHERE id = OLD.entry_id;

        WITH entry AS (
            SELECT user_id, category
            FROM mindwell.entries
            WHERE id = OLD.entry_id
        )
        UPDATE mindwell.vote_weights
        SET vote_count = vote_count - 1,
            vote_sum = vote_sum - OLD.vote,
            weight_sum = weight_sum - abs(OLD.vote),
            weight = CASE WHEN weight_sum = abs(OLD.vote) THEN 0.1
                ELSE atan2(vote_count - 1, 200) * (vote_sum - OLD.vote)
                    / (weight_sum - abs(OLD.vote)) / pi() * 2
                END
        FROM entry
        WHERE vote_weights.user_id = entry.user_id
            AND vote_weights.category = entry.category;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_entry_votes_ins
    AFTER INSERT ON mindwell.entry_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.entry_votes_ins();

CREATE TRIGGER cnt_entry_votes_upd
    AFTER UPDATE ON mindwell.entry_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.entry_votes_upd();

CREATE TRIGGER cnt_entry_votes_del
    AFTER DELETE ON mindwell.entry_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.entry_votes_del();


CREATE TABLE "mindwell"."vote_weights" (
	"user_id" Integer NOT NULL,
	"category" Integer NOT NULL,
    "weight" Real DEFAULT 0.1 NOT NULL,
    "vote_count" Integer DEFAULT 0 NOT NULL,
    "vote_sum" Real DEFAULT 0 NOT NULL,
    "weight_sum" Real DEFAULT 0 NOT NULL,
    CONSTRAINT "vote_weights_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "vote_weights_category" FOREIGN KEY("category") REFERENCES "mindwell"."categories"("id"),
    CONSTRAINT "unique_vote_weight" UNIQUE("user_id", "category") );
 ;

CREATE INDEX "index_vote_weights" ON "mindwell"."vote_weights" USING btree( "user_id" );

CREATE OR REPLACE FUNCTION mindwell.create_vote_weights() RETURNS TRIGGER AS $$
    BEGIN
        INSERT INTO mindwell.vote_weights(user_id, category) VALUES(NEW.id, 0);
        INSERT INTO mindwell.vote_weights(user_id, category) VALUES(NEW.id, 1);
        INSERT INTO mindwell.vote_weights(user_id, category) VALUES(NEW.id, 2);
        INSERT INTO mindwell.vote_weights(user_id, category) VALUES(NEW.id, 3);

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER crt_vote_weights
    AFTER INSERT ON mindwell.users
    FOR EACH ROW EXECUTE PROCEDURE mindwell.create_vote_weights();


CREATE TABLE "mindwell"."entries_privacy" (
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    CONSTRAINT "entries_privacy_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "entries_privacy_entry_id" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_entry_privacy" UNIQUE("user_id", "entry_id") );
 ;

CREATE INDEX "index_private_entries" ON "mindwell"."entries_privacy" USING btree( "entry_id" );

CREATE INDEX "index_private_users" ON "mindwell"."entries_privacy" USING btree( "user_id" );


CREATE TABLE "mindwell"."comments" (
	"id" Serial NOT NULL,
	"author_id" Integer NOT NULL,
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "edit_content" Text DEFAULT '' NOT NULL,
	"rating" Real DEFAULT 0 NOT NULL,
    "up_votes" Integer DEFAULT 0 NOT NULL,
    "down_votes" Integer DEFAULT 0 NOT NULL,
    "vote_sum" Real DEFAULT 0 NOT NULL,
    "weight_sum" Real DEFAULT 0 NOT NULL,
	CONSTRAINT "unique_comment_id" PRIMARY KEY( "id" ),
    CONSTRAINT "comment_author_id" FOREIGN KEY("author_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "comment_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "comment_entry_id" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE );
 ;

CREATE INDEX "index_comment_entry_id" ON "mindwell"."comments" USING btree( "entry_id" );

CREATE INDEX "index_comment_date" ON "mindwell"."comments" USING btree( "created_at" );


ALTER TABLE "mindwell"."entries"
ADD CONSTRAINT "entry_last_comment_id" FOREIGN KEY("last_comment") REFERENCES "mindwell"."comments"("id");

CREATE OR REPLACE FUNCTION mindwell.inc_comments() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET comments_count = comments_count + 1
        WHERE id = NEW.author_id;

        UPDATE mindwell.entries
        SET comments_count = comments_count + 1
        WHERE id = NEW.entry_id;

        INSERT INTO mindwell.watching
        VALUES(NEW.author_id, NEW.entry_id)
        ON CONFLICT ON CONSTRAINT unique_user_watching DO NOTHING;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comments_inc
    AFTER INSERT ON mindwell.comments
    FOR EACH ROW EXECUTE PROCEDURE mindwell.inc_comments();

CREATE OR REPLACE FUNCTION mindwell.dec_comments() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET comments_count = comments_count - 1
        WHERE id = OLD.author_id;

        UPDATE mindwell.entries
        SET comments_count = comments_count - 1
        WHERE id = OLD.entry_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comments_dec
    AFTER DELETE ON mindwell.comments
    FOR EACH ROW EXECUTE PROCEDURE mindwell.dec_comments();

CREATE OR REPLACE FUNCTION mindwell.set_last_comment_ins() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.entries
        SET last_comment = NEW.id
        WHERE id = NEW.entry_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER last_comments_ins
    AFTER INSERT ON mindwell.comments
    FOR EACH ROW EXECUTE PROCEDURE mindwell.set_last_comment_ins();

CREATE OR REPLACE FUNCTION mindwell.set_last_comment_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.entries
        SET last_comment = (
            SELECT max(comments.id)
            FROM comments
            WHERE entry_id = OLD.entry_id AND id <> OLD.id
        )
        WHERE last_comment = OLD.id;

        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER last_comments_del
    BEFORE DELETE ON mindwell.comments
    FOR EACH ROW EXECUTE PROCEDURE mindwell.set_last_comment_del();


CREATE TABLE "mindwell"."comment_votes" (
	"user_id" Integer NOT NULL,
	"comment_id" Integer NOT NULL,
    "vote" Real NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "comment_vote_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "comment_vote_comment_id" FOREIGN KEY("comment_id") REFERENCES "mindwell"."comments"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_comment_vote" UNIQUE("user_id", "comment_id") );
 ;

CREATE INDEX "index_voted_comments" ON "mindwell"."comment_votes" USING btree( "comment_id" );

CREATE INDEX "index_comment_voted_users" ON "mindwell"."comment_votes" USING btree( "user_id" );

CREATE OR REPLACE FUNCTION mindwell.comment_votes_ins() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.comments
        SET up_votes = up_votes + (NEW.vote > 0)::int,
            down_votes = down_votes + (NEW.vote < 0)::int,
            vote_sum = vote_sum + NEW.vote,
            weight_sum = weight_sum + abs(NEW.vote),
            rating = atan2(weight_sum + abs(NEW.vote), 2)
                * (vote_sum + NEW.vote) / (weight_sum + abs(NEW.vote)) / pi() * 200
        WHERE id = NEW.comment_id;

        WITH cmnt AS (
            SELECT user_id
            FROM mindwell.comments
            WHERE id = NEW.comment_id
        )
        UPDATE mindwell.vote_weights
        SET vote_count = vote_count + 1,
            vote_sum = vote_sum + NEW.vote,
            weight_sum = weight_sum + abs(NEW.vote),
            weight = atan2(vote_count + 1, 200) * (vote_sum + NEW.vote)
                / (weight_sum + abs(NEW.vote)) / pi() * 2
        FROM cmnt
        WHERE vote_weights.user_id = cmnt.user_id
            AND vote_weights.category =
                (SELECT id FROM categories WHERE "type" = 'comment');

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.comment_votes_upd() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.comments
        SET up_votes = up_votes - (OLD.vote > 0)::int + (NEW.vote > 0)::int,
            down_votes = down_votes - (OLD.vote < 0)::int + (NEW.vote < 0)::int,
            vote_sum = vote_sum - OLD.vote + NEW.vote,
            weight_sum = weight_sum - abs(OLD.vote) + abs(NEW.vote),
            rating = atan2(weight_sum - abs(OLD.vote) + abs(NEW.vote), 2)
                * (vote_sum - OLD.vote + NEW.vote) / (weight_sum - abs(OLD.vote) + abs(NEW.vote)) / pi() * 200
        WHERE id = NEW.comment_id;

        WITH cmnt AS (
            SELECT user_id
            FROM mindwell.comments
            WHERE id = NEW.comment_id
        )
        UPDATE mindwell.vote_weights
        SET vote_sum = vote_sum - OLD.vote + NEW.vote,
            weight_sum = weight_sum - abs(OLD.vote) + abs(NEW.vote),
            weight = atan2(vote_count, 200) * (vote_sum - OLD.vote + NEW.vote)
                / (weight_sum - abs(OLD.vote) + abs(NEW.vote)) / pi() * 2
        FROM cmnt
        WHERE vote_weights.user_id = cmnt.user_id
            AND vote_weights.category =
                (SELECT id FROM categories WHERE "type" = 'comment');

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.comment_votes_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.comments
        SET up_votes = up_votes - (OLD.vote > 0)::int,
            down_votes = down_votes - (OLD.vote < 0)::int,
            vote_sum = vote_sum - OLD.vote,
            weight_sum = weight_sum - abs(OLD.vote),
            rating = CASE WHEN weight_sum = abs(OLD.vote) THEN 0
                ELSE atan2(weight_sum - abs(OLD.vote), 2)
                    * (vote_sum - OLD.vote) / (weight_sum - abs(OLD.vote)) / pi() * 200
                END
        WHERE id = OLD.comment_id;

        WITH cmnt AS (
            SELECT user_id
            FROM mindwell.comments
            WHERE id = OLD.comment_id
        )
        UPDATE mindwell.vote_weights
        SET vote_count = vote_count - 1,
            vote_sum = vote_sum - OLD.vote,
            weight_sum = weight_sum - abs(OLD.vote),
            weight = CASE WHEN weight_sum = abs(OLD.vote) THEN 0.1
                ELSE atan2(vote_count - 1, 200) * (vote_sum - OLD.vote)
                    / (weight_sum - abs(OLD.vote)) / pi() * 2
                END
        FROM cmnt
        WHERE vote_weights.user_id = cmnt.user_id
            AND vote_weights.category =
                (SELECT id FROM categories WHERE "type" = 'comment');

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comment_votes_ins
    AFTER INSERT ON mindwell.comment_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.comment_votes_ins();

CREATE TRIGGER cnt_comment_votes_upd
    AFTER UPDATE ON mindwell.comment_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.comment_votes_upd();

CREATE TRIGGER cnt_comment_votes_del
    AFTER DELETE ON mindwell.comment_votes
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.comment_votes_del();


INSERT INTO mindwell.users
    (name, show_name, email, password_hash, invited_by, rank)
    VALUES('Mindwell', 'Mindwell', '', '', 1, 1);


CREATE TABLE "mindwell"."notification_type" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."notification_type" VALUES(0, 'comment');
INSERT INTO "mindwell"."notification_type" VALUES(1, 'follower');
INSERT INTO "mindwell"."notification_type" VALUES(2, 'request');
INSERT INTO "mindwell"."notification_type" VALUES(3, 'accept');
INSERT INTO "mindwell"."notification_type" VALUES(4, 'invite');
INSERT INTO "mindwell"."notification_type" VALUES(5, 'welcome');
INSERT INTO "mindwell"."notification_type" VALUES(6, 'invited');
INSERT INTO "mindwell"."notification_type" VALUES(7, 'adm_sent');
INSERT INTO "mindwell"."notification_type" VALUES(8, 'adm_received');
INSERT INTO "mindwell"."notification_type" VALUES(9, 'info');
INSERT INTO "mindwell"."notification_type" VALUES(10, 'wish_created');
INSERT INTO "mindwell"."notification_type" VALUES(11, 'wish_received');
INSERT INTO "mindwell"."notification_type" VALUES(12, 'entry_moved');
INSERT INTO "mindwell"."notification_type" VALUES(13, 'badge');


CREATE TABLE "mindwell"."notifications" (
    "id" Serial NOT NULL,
    "user_id" Integer NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "type" Integer NOT NULL,
    "subject_id" Integer NOT NULL,
    "read" Boolean DEFAULT FALSE NOT NULL,
	CONSTRAINT "unique_notification_id" PRIMARY KEY("id"),
    CONSTRAINT "notification_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_notification_type" FOREIGN KEY("type") REFERENCES "mindwell"."notification_type"("id") );
;

CREATE UNIQUE INDEX "index_notification_id" ON "mindwell"."notifications" USING btree( "id" );

CREATE INDEX "index_notification_user_id" ON "mindwell"."notifications" USING btree( "user_id" );


CREATE TABLE "mindwell"."info" (
    "id" Serial NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "content" Text NOT NULL,
    "link" Text NOT NULL,
	CONSTRAINT "unique_info_id" PRIMARY KEY("id") );
;


CREATE TABLE "mindwell"."images" (
	"id" Serial NOT NULL,
    "user_id" Integer NOT NULL,
	"path" Text NOT NULL,
    "extension" Text NOT NULL,
    "preview_extension" Text DEFAULT '' NOT NULL,
    "processing" Boolean DEFAULT TRUE NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT "unique_image_id" PRIMARY KEY("id"),
    CONSTRAINT "image_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"));
 ;

CREATE UNIQUE INDEX "index_image_id" ON "mindwell"."images" USING btree( "id" );

CREATE UNIQUE INDEX "index_image_path" ON "mindwell"."images" USING btree( "path" );


CREATE TABLE "mindwell"."size" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."size" VALUES(0, 'small');
INSERT INTO "mindwell"."size" VALUES(1, 'medium');
INSERT INTO "mindwell"."size" VALUES(2, 'large');
INSERT INTO "mindwell"."size" VALUES(3, 'thumbnail');


CREATE TABLE "mindwell"."image_sizes" (
    "image_id" Integer NOT NULL,
    "size" Integer NOT NULL,
    "width" Integer NOT NULL,
    "height" Integer NOT NULL,
    CONSTRAINT "unique_image_size" PRIMARY KEY ("image_id", "size"),
    CONSTRAINT "unique_image_id" FOREIGN KEY ("image_id") REFERENCES "mindwell"."images"("id") ON DELETE CASCADE,
    CONSTRAINT "enum_image_size" FOREIGN KEY("size") REFERENCES "mindwell"."size"("id")
);

CREATE INDEX "index_image_size_id" ON "mindwell"."image_sizes" USING btree( "image_id" );


CREATE TABLE "mindwell"."entry_images" (
    "entry_id" Integer NOT NULL,
    "image_id" Integer NOT NULL,
    "image_order" INTEGER NOT NULL,
    CONSTRAINT "entry_images_entry" FOREIGN KEY("entry_id") REFERENCES "mindwell"."entries"("id") ON DELETE CASCADE,
    CONSTRAINT "entry_images_image" FOREIGN KEY("image_id") REFERENCES "mindwell"."images"("id") ON DELETE CASCADE,
    CONSTRAINT "unique_entry_image" UNIQUE("entry_id", "image_id") );
;

CREATE INDEX "index_entry_images_entry" ON "mindwell"."entry_images" USING btree( "entry_id" );

CREATE INDEX "index_entry_images_image" ON "mindwell"."entry_images" USING btree( "image_id" );


CREATE TABLE "mindwell"."complain_type" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."complain_type" VALUES(0, 'comment');
INSERT INTO "mindwell"."complain_type" VALUES(1, 'entry');
INSERT INTO "mindwell"."complain_type" VALUES(2, 'message');
INSERT INTO "mindwell"."complain_type" VALUES(3, 'user');
INSERT INTO "mindwell"."complain_type" VALUES(4, 'theme');
INSERT INTO "mindwell"."complain_type" VALUES(5, 'wish');


CREATE TABLE "mindwell"."complains" (
    "id" Serial NOT NULL,
    "user_id" Integer NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "type" Integer NOT NULL,
    "subject_id" Integer NOT NULL,
    "content" Text NOT NULL,
    CONSTRAINT "unique_complain_id" PRIMARY KEY("id"),
    CONSTRAINT "complain_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_complain_type" FOREIGN KEY("type") REFERENCES "mindwell"."complain_type"("id") );
;

CREATE INDEX "index_complain_user_id" ON "mindwell"."complains" USING btree( "user_id" );


CREATE TABLE "mindwell"."chats" (
	"id" Serial NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"creator_id" Integer NOT NULL,
	"partner_id" Integer NOT NULL,
	"last_message" Integer,
	CONSTRAINT "unique_chat_id" PRIMARY KEY( "id" ),
	CONSTRAINT "chat_creator" FOREIGN KEY ("creator_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE,
	CONSTRAINT "chat_partner" FOREIGN KEY ("partner_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE,
	CONSTRAINT "unique_chat_partners" UNIQUE ( "creator_id", "partner_id" ) );
 ;

CREATE INDEX "index_chat_id" ON "mindwell"."chats" USING btree( "id" );

CREATE INDEX "index_chat_creator_id" ON "mindwell"."chats" USING btree( "creator_id" );

CREATE INDEX "index_chat_partner_id" ON "mindwell"."chats" USING btree( "partner_id" );

CREATE INDEX "index_last_message_id" ON "mindwell"."chats" USING btree( "last_message" );


CREATE TABLE "mindwell"."messages" (
	"id" Serial NOT NULL,
	"chat_id" Integer NOT NULL,
	"author_id" Integer NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"content" Text NOT NULL,
    "edit_content" Text NOT NULL,
	CONSTRAINT "unique_message_id" PRIMARY KEY( "id" ),
    CONSTRAINT "message_user_id" FOREIGN KEY("author_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE,
    CONSTRAINT "message_chat_id" FOREIGN KEY("chat_id") REFERENCES "mindwell"."chats"("id") ON DELETE CASCADE );
 ;

CREATE INDEX "index_message_id" ON "mindwell"."messages" USING btree( "id" );

CREATE INDEX "index_message_chat" ON "mindwell"."messages" USING btree( "chat_id" );


ALTER TABLE "mindwell"."chats"
ADD CONSTRAINT "chat_last_message_id" FOREIGN KEY("last_message") REFERENCES "mindwell"."messages"("id");

CREATE OR REPLACE FUNCTION mindwell.set_last_message_ins() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.chats
        SET last_message = NEW.id
        WHERE id = NEW.chat_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER last_messages_ins
    AFTER INSERT ON mindwell.messages
    FOR EACH ROW EXECUTE PROCEDURE mindwell.set_last_message_ins();

CREATE OR REPLACE FUNCTION mindwell.set_last_message_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.chats
        SET last_message = (
            SELECT max(messages.id)
            FROM messages
            WHERE chat_id = OLD.chat_id AND id <> OLD.id
        )
        WHERE last_message = OLD.id;

        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER last_messages_del
    BEFORE DELETE ON mindwell.messages
    FOR EACH ROW EXECUTE PROCEDURE mindwell.set_last_message_del();


CREATE TABLE "mindwell"."talkers" (
	"chat_id" Integer NOT NULL,
	"user_id" Integer NOT NULL,
	"last_read" Integer DEFAULT 0 NOT NULL,
	"unread_count" Integer DEFAULT 0 NOT NULL,
	CONSTRAINT "unique_talker_chat" PRIMARY KEY( "chat_id", "user_id" ),
    CONSTRAINT "talkers_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE,
    CONSTRAINT "talkers_chat_id" FOREIGN KEY("chat_id") REFERENCES "mindwell"."chats"("id") ON DELETE CASCADE);
 ;

CREATE INDEX "index_talkers_chat" ON "mindwell"."talkers" USING btree( "chat_id" );

CREATE INDEX "index_talkers_user" ON "mindwell"."talkers" USING btree( "user_id" );

CREATE OR REPLACE FUNCTION mindwell.count_unread_ins() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.talkers
        SET unread_count = unread_count + 1
        WHERE talkers.chat_id = NEW.chat_id AND talkers.user_id <> NEW.author_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_unread_ins
    AFTER INSERT ON mindwell.messages
    FOR EACH ROW EXECUTE PROCEDURE mindwell.count_unread_ins();

CREATE OR REPLACE FUNCTION mindwell.count_unread_del() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.talkers
        SET unread_count = unread_count - 1
        WHERE talkers.chat_id = OLD.chat_id AND talkers.user_id <> OLD.author_id
            AND last_read < OLD.id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_unread_del
    AFTER DELETE ON mindwell.messages
    FOR EACH ROW EXECUTE PROCEDURE mindwell.count_unread_del();


CREATE TABLE "mindwell"."badges" (
    "id" Serial NOT NULL,
    "code" Text NOT NULL,
    "level" Integer DEFAULT 1 NOT NULL,
    "title" Text NOT NULL,
    "description" Text NOT NULL,
    "icon" Text NOT NULL,
    "user_count" Integer DEFAULT 0 NOT NULL,
    CONSTRAINT "unique_badge_id" PRIMARY KEY ( "id" ),
    CONSTRAINT "unique_badge_code_level" UNIQUE ( "code", "level" )
);

INSERT INTO "mindwell"."badges" ("code", "title", "description", "icon")
VALUES ('test1', 'test1', 'test badge 1', 'test1.webp'),
       ('test2', 'test2', 'test badge 2', 'test2.webp');


CREATE TABLE "mindwell"."user_badges" (
    "user_id" Integer NOT NULL,
    "badge_id" Integer NOT NULL,
    "given_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "unique_user_badge_badge_id" FOREIGN KEY ("badge_id") REFERENCES mindwell.badges( "id" ),
    CONSTRAINT "unique_user_badge_user_id" FOREIGN KEY ( "user_id" ) REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "unique_user_badge" UNIQUE ( "user_id", "badge_id" )
);

CREATE OR REPLACE FUNCTION mindwell.inc_badges() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET badges_count = badges_count + 1
        WHERE id = NEW.user_id;

        UPDATE mindwell.badges
        SET user_count = user_count + 1
        WHERE id = NEW.badge_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.dec_badges() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET badges_count = badges_count - 1
        WHERE id = OLD.user_id;

        UPDATE mindwell.badges
        SET user_count = user_count - 1
        WHERE id = OLD.badge_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_badges_inc
    AFTER INSERT ON mindwell.user_badges
    FOR EACH ROW EXECUTE PROCEDURE mindwell.inc_badges();

CREATE TRIGGER cnt_badges_dec
    AFTER DELETE ON mindwell.user_badges
    FOR EACH ROW EXECUTE PROCEDURE mindwell.dec_badges();

CREATE OR REPLACE FUNCTION mindwell.notify_user_badges() RETURNS TRIGGER AS $$
    BEGIN
            PERFORM pg_notify('user_badges', json_build_object(
                'user_id', NEW.user_id,
                'badge_id', NEW.badge_id
            )::text);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ntf_badges
    AFTER INSERT ON mindwell.user_badges
    FOR EACH ROW EXECUTE FUNCTION mindwell.notify_user_badges();


CREATE TABLE "mindwell"."apps" (
    "id" Integer UNIQUE NOT NULL,
    "secret_hash" Bytea NOT NULL,
    "redirect_uri" Text NOT NULL,
    "developer_id" Integer NOT NULL,
    "flow" Smallint NOT NULL,
    "name" Text NOT NULL,
    "show_name" Text NOT NULL,
    "platform" Text NOT NULL,
    "info" Text NOT NULL,
    "ban" Bool NOT NULL DEFAULT FALSE,
     CONSTRAINT "unique_app_id" PRIMARY KEY("id"),
     CONSTRAINT "app_developer_id" FOREIGN KEY("developer_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE
);


CREATE TABLE "mindwell"."sessions" (
    "id" Bigserial,
    "app_id" Integer NOT NULL,
    "user_id" Integer NOT NULL,
    "scope" Integer NOT NULL,
    "access_hash" Bytea NOT NULL,
    "refresh_hash" Bytea NOT NULL,
    "access_thru" Timestamp With Time Zone NOT NULL,
    "refresh_thru" Timestamp With Time Zone NOT NULL,
    CONSTRAINT "session_id" PRIMARY KEY("id"),
    CONSTRAINT "session_app_id" FOREIGN KEY("app_id") REFERENCES "mindwell"."apps"("id") ON DELETE CASCADE,
    CONSTRAINT "session_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id") ON DELETE CASCADE
);

CREATE INDEX "index_access_hash" ON "mindwell"."sessions" USING btree( "access_hash" );
CREATE INDEX "index_refresh_hash" ON "mindwell"."sessions" USING btree( "refresh_hash" );


CREATE TABLE "mindwell"."app_tokens" (
    "app_id" Integer NOT NULL,
    "token_hash" Bytea NOT NULL,
    "valid_thru" Timestamp With Time Zone NOT NULL,
     CONSTRAINT "app_token_app_id" FOREIGN KEY("app_id") REFERENCES "mindwell"."apps"("id") ON DELETE CASCADE
);

CREATE INDEX "index_app_token_hash" ON "mindwell"."app_tokens" USING btree( "token_hash" );


CREATE TABLE "mindwell"."user_log" (
    "name" Text NOT NULL,
    "user_agent" Text NOT NULL,
    "ip" Inet NOT NULL,
    "device" Integer NOT NULL,
    "app" Bigint NOT NULL,
    "uid" Integer NOT NULL,
    "uid2" Bigint NOT NULL,
    "at" Timestamp With Time Zone NOT NULL,
    "first" Boolean NOT NULL
);

CREATE VIEW "mindwell"."user_log_view" AS
SELECT name, ip, to_hex(device) AS device, to_hex(app) AS app, to_hex(uid) AS uid, to_hex(uid2) AS uid2,
       to_char(at, 'YYYY.MM.DD HH24:MI:SS') AS at, first AS f
FROM mindwell.user_log
ORDER BY at DESC;

CREATE INDEX "index_requested_at" ON "mindwell"."user_log" USING btree( "at" );

CREATE INDEX "index_user_log_name" ON "mindwell"."user_log" USING btree( "name" );


CREATE OR REPLACE FUNCTION mindwell.give_invites() RETURNS TABLE(user_id int) AS $$
WITH inviters AS (
    UPDATE mindwell.users
        SET last_invite = CURRENT_DATE
        WHERE ((id IN (
            SELECT user_id
            FROM (
                     SELECT entries.created_at, user_id
                     FROM mindwell.entries
                              JOIN mindwell.users ON user_id = users.id
                     WHERE age(entries.created_at) <= interval '1 month'
                       AND visible_for in (
                         SELECT id
                         FROM mindwell.entry_privacy
                         WHERE type in ('all', 'registered', 'invited')
                     )
                       AND users.invited_by IS NOT NULL
                       AND users.privacy in (
                         SELECT id
                         FROM mindwell.user_privacy
                         WHERE type in ('all', 'registered', 'invited')
                     )
                     ORDER BY rating DESC
                     LIMIT 100) AS e
            WHERE current_timestamp - e.created_at < interval '3 days'
        )
            AND (SELECT COUNT(*) FROM mindwell.invites WHERE referrer_id = users.id) < 3
                   ) OR (
                           last_invite = created_at::Date
                       AND (
                               SELECT COUNT(DISTINCT entries.id)
                               FROM mindwell.entries
                                        JOIN mindwell.entry_votes ON entries.id = entry_votes.entry_id
                               WHERE entries.user_id = users.id
                                 AND entry_votes.vote > 0
                                 AND entry_votes.user_id <> users.invited_by
                           ) >= 10
                   )) AND age(last_invite) >= interval '14 days'
            AND invite_ban <= CURRENT_DATE
        RETURNING users.id
), wc AS (
    SELECT COUNT(*) AS words FROM mindwell.invite_words
)
INSERT INTO mindwell.invites(referrer_id, word1, word2, word3)
SELECT inviters.id,
       trunc(random() * wc.words),
       trunc(random() * wc.words),
       trunc(random() * wc.words)
FROM inviters, wc
ON CONFLICT (word1, word2, word3) DO NOTHING
RETURNING referrer_id;
$$ LANGUAGE SQL;


CREATE OR REPLACE FUNCTION mindwell.recalc_karma() RETURNS VOID AS $$
    WITH upd AS (
        SELECT users.id, (
                users.karma * 4
                + COALESCE(fek.karma, 0) + COALESCE(bek.karma, 0)
                + COALESCE(fck.karma, 0) / 10 + COALESCE(bck.karma, 0) / 10
            ) / 5 AS karma
        FROM mindwell.users
        LEFT JOIN (
            SELECT entries.user_id AS id, sum(entry_votes.vote) AS karma
            FROM mindwell.entries
            JOIN mindwell.entry_votes ON entry_votes.entry_id = entries.id
            WHERE abs(entry_votes.vote) > 0.2 AND age(entries.created_at) <= interval '2 months'
            GROUP BY entries.user_id
        ) AS fek ON users.id = fek.id -- votes for users entries
        LEFT JOIN (
            SELECT entry_votes.user_id AS id, sum(entry_votes.vote) / 5 AS karma
            FROM mindwell.entry_votes
            WHERE entry_votes.vote < 0 AND age(entry_votes.created_at) <= interval '2 months'
            GROUP BY entry_votes.user_id
        ) AS bek ON users.id = bek.id -- entry votes by users
        LEFT JOIN (
            SELECT comments.user_id AS id, sum(comment_votes.vote) AS karma
            FROM mindwell.comments
            JOIN mindwell.comment_votes ON comment_votes.comment_id = comments.id
            WHERE abs(comment_votes.vote) > 0.2 AND age(comments.created_at) <= interval '2 months'
            GROUP BY comments.user_id
        ) AS fck ON users.id = fck.id -- votes for users comments
        LEFT JOIN (
            SELECT comment_votes.user_id AS id, sum(comment_votes.vote) / 5 AS karma
            FROM mindwell.comment_votes
            WHERE comment_votes.vote < 0 AND age(comment_votes.created_at) <= interval '2 months'
            GROUP BY comment_votes.user_id
        ) AS bck ON users.id = bck.id -- comment votes by users
        WHERE users.creator_id IS NULL
    )
    UPDATE mindwell.users
    SET karma = upd.karma
    FROM upd
    WHERE users.id = upd.id AND users.karma <> upd.karma;

    WITH upd AS (
        SELECT id, row_number() OVER (ORDER BY karma DESC, created_at ASC) as rank
        FROM mindwell.users
        WHERE users.creator_id IS NULL
    )
    UPDATE mindwell.users
    SET rank = upd.rank
    FROM upd
    WHERE users.id = upd.id AND users.rank <> upd.rank;

    WITH upd AS (
        SELECT users.id, (
                users.karma * 4
                + COALESCE(fek.karma, 0)
                + COALESCE(fck.karma, 0) / 10
            ) / 5 AS karma
        FROM mindwell.users
        LEFT JOIN (
            SELECT entries.author_id AS id, sum(entry_votes.vote) AS karma
            FROM mindwell.entries
            JOIN mindwell.entry_votes ON entry_votes.entry_id = entries.id
            WHERE entries.author_id <> entries.user_id
                AND entry_votes.vote > 0.2 AND age(entries.created_at) <= interval '2 months'
            GROUP BY entries.author_id
        ) AS fek ON users.id = fek.id -- votes for users entries
        LEFT JOIN (
            SELECT comments.author_id AS id, sum(comment_votes.vote) AS karma
            FROM mindwell.comments
            JOIN mindwell.comment_votes ON comment_votes.comment_id = comments.id
            WHERE comments.author_id <> comments.user_id
                AND comment_votes.vote > 0.2 AND age(comments.created_at) <= interval '2 months'
            GROUP BY comments.author_id
        ) AS fck ON users.id = fck.id -- votes for users comments
        WHERE users.creator_id IS NOT NULL
    )
    UPDATE mindwell.users
    SET karma = upd.karma
    FROM upd
    WHERE users.id = upd.id AND users.karma <> upd.karma;

    WITH upd AS (
        SELECT id, row_number() OVER (ORDER BY karma DESC, created_at ASC) as rank
        FROM mindwell.users
        WHERE users.creator_id IS NOT NULL
    )
    UPDATE mindwell.users
    SET rank = upd.rank
    FROM upd
    WHERE users.id = upd.id AND users.rank <> upd.rank;
$$ LANGUAGE SQL;


CREATE OR REPLACE FUNCTION mindwell.ban_invite(userName Text) RETURNS VOID AS $$
    DELETE FROM mindwell.invites
    WHERE referrer_id = (SELECT id FROM mindwell.users WHERE lower(name) = lower(userName));

    UPDATE mindwell.users
    SET invite_ban = CURRENT_DATE + interval '1 month'
    WHERE lower(name) = lower(userName);
$$ LANGUAGE SQL;


CREATE OR REPLACE FUNCTION mindwell.delete_user(user_name TEXT) RETURNS VOID AS $$
    DECLARE
        del_id INTEGER;
    BEGIN
        del_id = (SELECT id FROM users WHERE lower(name) = lower(user_name));

        DELETE FROM mindwell.relations WHERE to_id = del_id;
        DELETE FROM mindwell.relations WHERE from_id = del_id;

        DELETE FROM mindwell.favorites WHERE favorites.user_id = del_id;
        DELETE FROM mindwell.watching WHERE watching.user_id = del_id;
        DELETE FROM mindwell.entries_privacy WHERE entries_privacy.user_id = del_id;

        DELETE FROM mindwell.entry_votes WHERE entry_votes.user_id = del_id;
        DELETE FROM mindwell.comment_votes WHERE comment_votes.user_id = del_id;
        DELETE FROM mindwell.vote_weights WHERE vote_weights.user_id = del_id;

        DELETE FROM mindwell.notifications
        WHERE notifications.user_id = del_id OR
            CASE (SELECT "type" FROM notification_type WHERE notification_type.id = notifications."type")
            WHEN 'comment' THEN
                (SELECT user_id FROM comments WHERE comments.id = notifications.subject_id) = del_id
            WHEN 'invite' THEN
                FALSE
            ELSE
                notifications.subject_id = del_id
            END;

        DELETE FROM complains WHERE user_id = del_id;

        DELETE FROM mindwell.images WHERE images.user_id = del_id;
        DELETE FROM mindwell.entries WHERE user_id = del_id;
        DELETE FROM mindwell.comments WHERE user_id = del_id;
        DELETE FROM mindwell.users WHERE id = del_id;
    END;
$$ LANGUAGE plpgsql;
