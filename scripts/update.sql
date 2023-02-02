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
