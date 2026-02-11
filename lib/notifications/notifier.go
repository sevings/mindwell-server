package notifications

import (
	"context"
	"encoding/json"
	"log"

	"github.com/centrifugal/gocent"
	"github.com/sevings/mindwell-server/lib/database"
)

type message struct {
	ID    int64  `json:"id"`
	Subj  int64  `json:"subject,omitempty"`
	Type  string `json:"type,omitempty"`
	State string `json:"state,omitempty"`
	ch    string
}

type Notifier struct {
	cent *gocent.Client
	ch   chan *message
	stop chan any
}

const (
	stateNew     = "new"
	stateUpdated = "updated"
	stateRemoved = "removed"
	stateRead    = "read"

	TypeAdmReceived  = "adm_received"
	TypeAdmSent      = "adm_sent"
	TypeInvited      = "invited"
	TypeAccept       = "accept"
	TypeComment      = "comment"
	TypeRequest      = "request"
	TypeFollower     = "follower"
	TypeInvite       = "invite"
	TypeMessage      = "message"
	TypeWishReceived = "wish_received"
	TypeWishCreated  = "wish_created"
	TypeEntryMoved   = "entry_moved"
	TypeBadge        = "badge"
)

func NewNotifier(apiURL, apiKey string) *Notifier {
	if len(apiKey) == 0 {
		return &Notifier{}
	}

	cfg := gocent.Config{
		Addr: apiURL,
		Key:  apiKey,
	}

	ntf := &Notifier{
		cent: gocent.New(cfg),
		ch:   make(chan *message, 200),
		stop: make(chan any),
	}

	go func() {
		ctx := context.Background()

		for msg := range ntf.ch {
			data, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				continue
			}

			err = ntf.cent.Publish(ctx, msg.ch, data)
			if err != nil {
				log.Println(err)
			}
		}

		close(ntf.stop)
	}()

	return ntf
}

func (ntf *Notifier) Stop() {
	if ntf.ch == nil {
		return
	}

	close(ntf.ch)
	<-ntf.stop
}

func notificationsChannel(userName string) string {
	return "notifications#" + userName
}

func (ntf *Notifier) Notify(tx *database.AutoTx, subjectID int64, tpe, user string) {
	const q = `
		INSERT INTO notifications(user_id, subject_id, type)
		VALUES((SELECT id from users WHERE lower(name) = lower($1)),
			$2, (SELECT id FROM notification_type WHERE type = $3))
		RETURNING id
	`

	var id int64
	tx.Query(q, user, subjectID, tpe).Scan(&id)

	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    id,
			Subj:  subjectID,
			State: stateNew,
			Type:  tpe,
			ch:    notificationsChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyUpdate(tx *database.AutoTx, subjectID int64, tpe string) {
	if ntf.ch == nil {
		return
	}

	const q = `
		SELECT notifications.id, name
		FROM notifications, users
		WHERE subject_id = $1 AND type = (SELECT id FROM notification_type WHERE type = $2)
			AND user_id = users.id
	`

	tx.Query(q, subjectID, tpe)

	for {
		var id int64
		var user string
		ok := tx.Scan(&id, &user)
		if !ok {
			break
		}

		ntf.ch <- &message{
			ID:    id,
			State: stateUpdated,
			ch:    notificationsChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyRemove(tx *database.AutoTx, subjectID int64, tpe string) {
	const q = `
		DELETE FROM notifications
		WHERE subject_id = $1 AND type = (SELECT id FROM notification_type WHERE type = $2)
		RETURNING id, (SELECT name FROM users WHERE id = user_id)
	`

	tx.Query(q, subjectID, tpe)

	for {
		var id int64
		var user string
		ok := tx.Scan(&id, &user)
		if !ok {
			break
		}

		if ntf.ch == nil {
			continue
		}

		ntf.ch <- &message{
			ID:    id,
			State: stateRemoved,
			ch:    notificationsChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyRead(user string, ntfID int64) {
	if ntf.ch == nil {
		return
	}

	ntf.ch <- &message{
		ID:    ntfID,
		State: stateRead,
		ch:    notificationsChannel(user),
	}
}

func (ntf *Notifier) NotifyNewFollower(tx *database.AutoTx, fromID int64, to string, isPrivate bool) {
	tpe := TypeFollower
	if isPrivate {
		tpe = TypeRequest
	}

	ntf.Notify(tx, fromID, tpe, to)
}

func (ntf *Notifier) NotifyRemoveFollower(tx *database.AutoTx, fromID, toID int64, to string) {
	const q = `
		DELETE FROM notifications
	    WHERE id = (SELECT MAX(id) FROM notifications
			WHERE subject_id = $1
			  AND user_id = $2
			  AND type IN (SELECT id FROM notification_type WHERE type = $3 OR type = $4)
		)
		RETURNING id
	`

	ntfID := tx.QueryInt64(q, fromID, toID, TypeFollower, TypeRequest)

	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    ntfID,
			State: stateRemoved,
			ch:    notificationsChannel(to),
		}
	}
}

func messagesChannel(userName string) string {
	return "messages#" + userName
}

func (ntf *Notifier) NotifyMessage(chatID, msgID int64, user string) {
	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    chatID,
			Subj:  msgID,
			State: stateNew,
			Type:  TypeMessage,
			ch:    messagesChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyMessageUpdate(chatID, msgID int64, user string) {
	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    chatID,
			Subj:  msgID,
			State: stateUpdated,
			Type:  TypeMessage,
			ch:    messagesChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyMessageRemove(chatID, msgID int64, user string) {
	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    chatID,
			Subj:  msgID,
			State: stateRemoved,
			Type:  TypeMessage,
			ch:    messagesChannel(user),
		}
	}
}

func (ntf *Notifier) NotifyMessageRead(chatID, msgID int64, user string) {
	if ntf.ch != nil {
		ntf.ch <- &message{
			ID:    chatID,
			Subj:  msgID,
			State: stateRead,
			Type:  TypeMessage,
			ch:    messagesChannel(user),
		}
	}
}
