package wishes

import (
	"database/sql"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
	"time"
)

func LastCreatedWish(db *sql.DB, userID *models.UserID) (int64, bool) {
	tx := utils.NewAutoTx(db)
	defer tx.Finish()

	const q = `
SELECT wishes.id
FROM wishes
JOIN wish_states ON wishes.state = wish_states.id
WHERE from_id = $1
	AND wish_states.state = 'new' AND NOW() - created_at < interval '8 hours'
ORDER BY created_at DESC
LIMIT 1
`

	wishID := tx.QueryInt64(q, userID.ID)
	return wishID, wishID != 0
}

func LoadWish(tx *utils.AutoTx, userID *models.UserID, id int64) (*models.Wish, bool) {
	const q = `
SELECT content, to_id, wish_states.state,
	extract(epoch from created_at + interval '8 hours')
FROM wishes
JOIN wish_states ON wishes.state = wish_states.id
WHERE wishes.id = $1 AND (
    from_id = $2
	OR (to_id = $3 AND wish_states.state IN ('sent', 'thanked'))
)
`

	var receiverID int64
	wish := &models.Wish{
		ID: id,
	}
	tx.Query(q, id, userID.ID, userID.ID).
		Scan(&wish.Content, &receiverID, &wish.State,
			&wish.SendUntil)

	if tx.Error() != nil {
		return nil, false
	}

	if userID.ID == receiverID {
		wish.SendUntil = 0
	} else {
		wish.Receiver = utils.LoadUser(tx, receiverID)
		switch wish.State {
		case models.WishStateNew:
			if time.Now().Unix() > int64(wish.SendUntil) {
				wish.State = models.WishStateExpired
			}
		case models.WishStateComplained:
			fallthrough
		case models.WishStateThanked:
			wish.State = models.WishStateSent
			fallthrough
		case models.WishStateSent:
			fallthrough
		case models.WishStateDeclined:
			wish.SendUntil = 0
		}
	}

	return wish, true
}

func saveWish(tx *utils.AutoTx, userID *models.UserID, id int64, content string) (int64, bool) {
	const q = `
UPDATE wishes
SET content = $3,
    state = (SELECT id FROM wish_states WHERE state = 'sent')
WHERE wishes.id = $1 AND from_id = $2
	AND state = (SELECT id FROM wish_states WHERE state = 'new')
	AND NOW() - created_at < interval '8 hours'
RETURNING to_id
`

	toID := tx.QueryInt64(q, id, userID.ID, content)
	return toID, toID != 0
}

func declineWish(tx *utils.AutoTx, userID *models.UserID, id int64) bool {
	const q = `
UPDATE wishes
SET state = (SELECT id FROM wish_states WHERE state = 'declined')
WHERE wishes.id = $1 AND from_id = $2
	AND state IN (SELECT id FROM wish_states WHERE state IN ('new','declined'))
	AND NOW() - created_at < interval '8 hours'
`

	tx.Exec(q, id, userID.ID)
	return tx.RowsAffected() == 1
}

func thankWish(tx *utils.AutoTx, userID *models.UserID, id int64) bool {
	const q = `
UPDATE wishes
SET state = (SELECT id FROM wish_states WHERE state = 'thanked')
WHERE wishes.id = $1 AND to_id = $2
	AND state IN (SELECT id FROM wish_states WHERE state IN ('sent', 'thanked'))
`

	tx.Exec(q, id, userID.ID)
	return tx.RowsAffected() == 1
}

func createWish(tx *utils.AutoTx, userID int64) (int64, bool) {
	const canSendQuery = `
SELECT invited_by IS NOT NULL
	AND invite_ban <= CURRENT_DATE
	AND vote_ban <= CURRENT_DATE
	AND comment_ban <= CURRENT_DATE
	AND live_ban <= CURRENT_DATE
	AND complain_ban <= CURRENT_DATE
	AND user_ban <= CURRENT_DATE
	AND NOT adm_ban
	AND karma > 0
    AND verified
	AND send_wishes
	AND NOW() - COALESCE(prev_wish.created_at, users.created_at) > interval '3 days'
	AND users.authority = (SELECT id FROM authority WHERE type = 'user')
FROM users
LEFT JOIN (
    SELECT from_id, created_at
    FROM wishes
    WHERE from_id = $1
    ORDER BY id DESC
    LIMIT 1
) AS prev_wish ON users.id = prev_wish.from_id
WHERE users.id = $1
`

	if !tx.QueryBool(canSendQuery, userID) {
		return 0, false
	}

	const receiverIDQuery = `
SELECT users.id
FROM users
LEFT JOIN relations AS from_me ON users.id = from_me.to_id AND from_me.from_id = $1
LEFT JOIN relation AS rel_from_me ON from_me.type = rel_from_me.id
LEFT JOIN relations AS to_me ON users.id = to_me.from_id AND to_me.to_id = $1
LEFT JOIN relation AS rel_to_me ON to_me.type = rel_to_me.id
JOIN authority ON users.authority = authority.id
LEFT JOIN (
    SELECT to_id, MAX(created_at) created_at
    FROM wishes
	JOIN wish_states ON wishes.state = wish_states.id
    WHERE wish_states.state = 'sent' OR wish_states.state = 'thanked'
    GROUP BY to_id
) AS last_wishes ON users.id = last_wishes.to_id
LEFT JOIN (
    SELECT DISTINCT to_id
    FROM wishes
    WHERE from_id = $1
    	AND NOW() - wishes.created_at < interval '2 months'
) AS sent_wishes ON users.id = sent_wishes.to_id
WHERE NOW() - last_seen_at < interval '1 day'
	AND invited_by IS NOT NULL
	AND invite_ban <= CURRENT_DATE
	AND vote_ban <= CURRENT_DATE
	AND comment_ban <= CURRENT_DATE
	AND live_ban <= CURRENT_DATE
	AND complain_ban <= CURRENT_DATE
	AND user_ban <= CURRENT_DATE
	AND NOT adm_ban
	AND karma > 0
  	AND verified
	AND send_wishes
	AND COALESCE(rel_from_me.type, 'none') <> 'ignored'
	AND COALESCE(rel_to_me.type, 'none') <> 'ignored'
  	AND authority.type = 'user'
	AND sent_wishes.to_id IS NULL
	AND users.id <> $1
ORDER BY COALESCE(last_wishes.created_at, users.created_at)::DATE,
	COALESCE(rel_from_me.type, 'none') = 'followed' DESC
LIMIT 1
`

	receiverID := tx.QueryInt64(receiverIDQuery, userID)
	if receiverID == 0 {
		return 0, false
	}

	const createWishQuery = `
INSERT INTO wishes(from_id, to_id)
VALUES ($1, $2)
RETURNING id
`

	wishID := tx.QueryInt64(createWishQuery, userID, receiverID)

	return wishID, true
}
