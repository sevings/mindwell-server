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
