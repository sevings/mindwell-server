ALTER TABLE users
ADD COLUMN "creator_id" Integer;

ALTER TABLE users
ADD CONSTRAINT "theme_creator" FOREIGN KEY("creator_id") REFERENCES "mindwell"."users"("id");

ALTER TABLE entries
ADD COLUMN "user_id" Integer NOT NULL DEFAULT 1;

UPDATE entries
SET user_id = author_id;

ALTER TABLE entries
ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE entries
DROP CONSTRAINT entry_user_id;

ALTER TABLE entries
ADD CONSTRAINT "entry_author_id" FOREIGN KEY("author_id") REFERENCES "mindwell"."users"("id");

ALTER TABLE entries
ADD CONSTRAINT "entry_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id");

DROP INDEX index_entry_users_id;

CREATE INDEX "index_entry_author_id" ON "mindwell"."entries" USING btree( "author_id" );

CREATE INDEX "index_entry_user_id" ON "mindwell"."entries" USING btree( "user_id" );

ALTER TABLE entries
ADD COLUMN "is_anonymous" Boolean DEFAULT FALSE NOT NULL;

ALTER TABLE comments
ADD COLUMN "user_id" Integer NOT NULL DEFAULT 1;

UPDATE comments
SET user_id = author_id;

ALTER TABLE comments
ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE comments
DROP CONSTRAINT comment_user_id;

ALTER TABLE comments
ADD CONSTRAINT "comment_author_id" FOREIGN KEY("author_id") REFERENCES "mindwell"."users"("id");

ALTER TABLE comments
ADD CONSTRAINT "comment_user_id" FOREIGN KEY("user_id") REFERENCES "mindwell"."users"("id");

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
            weight = atan2(vote_count + 1, 20) * (vote_sum + NEW.vote)
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
            weight = atan2(vote_count, 20) * (vote_sum - OLD.vote + NEW.vote)
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
                ELSE atan2(vote_count - 1, 20) * (vote_sum - OLD.vote)
                    / (weight_sum - abs(OLD.vote)) / pi() * 2
                END
        FROM entry
        WHERE vote_weights.user_id = entry.user_id
            AND vote_weights.category = entry.category;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

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
            weight = atan2(vote_count + 1, 20) * (vote_sum + NEW.vote)
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
            weight = atan2(vote_count, 20) * (vote_sum - OLD.vote + NEW.vote)
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
                ELSE atan2(vote_count - 1, 20) * (vote_sum - OLD.vote)
                    / (weight_sum - abs(OLD.vote)) / pi() * 2
                END
        FROM cmnt
        WHERE vote_weights.user_id = cmnt.user_id
            AND vote_weights.category =
                (SELECT id FROM categories WHERE "type" = 'comment');

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

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
                AND abs(entry_votes.vote) > 0.2 AND age(entries.created_at) <= interval '2 months'
            GROUP BY entries.author_id
        ) AS fek ON users.id = fek.id -- votes for users entries
        LEFT JOIN (
            SELECT comments.author_id AS id, sum(comment_votes.vote) AS karma
            FROM mindwell.comments
            JOIN mindwell.comment_votes ON comment_votes.comment_id = comments.id
            WHERE comments.author_id <> comments.user_id
                AND abs(comment_votes.vote) > 0.2 AND age(comments.created_at) <= interval '2 months'
            GROUP BY comments.author_id
        ) AS fck ON users.id = fck.id -- votes for users comments
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
