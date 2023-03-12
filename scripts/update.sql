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
