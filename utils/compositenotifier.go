package utils

import (
	"database/sql"
	"encoding/json"

	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
	"gitlab.com/golang-commonmark/markdown"
)

type CompositeNotifier struct {
	srv  *MindwellServer
	md   *markdown.Markdown
	Mail MailSender
	Ntf  *Notifier
	Tg   *TelegramBot
}

func NewCompositeNotifier(srv *MindwellServer) *CompositeNotifier {
	ntfURL := srv.ConfigString("centrifugo.api_url")
	ntfKey := srv.ConfigString("centrifugo.api_key")

	ntf := &CompositeNotifier{
		srv: srv,
		md:  markdown.New(markdown.Typographer(false), markdown.Breaks(true), markdown.Tables(false)),
		Ntf: NewNotifier(ntfURL, ntfKey),
		Tg:  NewTelegramBot(srv),
	}

	srv.PS.Subscribe("moved_entries", ntf.notifyMovedEntry)
	srv.PS.Subscribe("user_badges", ntf.notifyBadge)

	return ntf
}

func (ntf *CompositeNotifier) Stop() {
	ntf.Mail.Stop()
	ntf.Tg.Stop()
	ntf.Ntf.Stop()
}

func (ntf *CompositeNotifier) SendNewInvite(tx *AutoTx, userID int) {
	var email, name, showName string
	var sendEmail, sendTg bool
	var tg sql.NullInt64

	const q = `
		SELECT email, name, show_name, verified AND email_invites, telegram, telegram_invites
		FROM users WHERE id = $1
	`

	tx.Query(q, userID)
	tx.Scan(&email, &name, &showName, &sendEmail, &tg, &sendTg)

	ntf.Ntf.Notify(tx, 0, typeInvite, name)

	if tg.Valid && sendTg {
		ntf.Tg.SendNewInvite(tg.Int64)
	}

	if sendEmail {
		ntf.Mail.SendNewInvite(email, showName)
	}
}

func (ntf *CompositeNotifier) SendEmailChanged(tx *AutoTx, userID *models.UserID, oldEmail, newEmail string) {
	const q = `
		SELECT show_name, telegram
		FROM users
		WHERE id = $1
	`

	var name string
	var tg sql.NullInt64
	tx.Query(q, userID.ID).Scan(&name, &tg)

	if tx.Error() != nil {
		if tx.Error() != sql.ErrNoRows {
			ntf.srv.LogSystem().Error(tx.Error().Error())
		}
		return
	}

	if len(oldEmail) > 0 {
		ntf.Mail.SendEmailChanged(oldEmail, name)
	}

	code := ntf.srv.TokenHash().VerificationCode(newEmail)
	ntf.Mail.SendGreeting(newEmail, name, code)

	if tg.Valid {
		ntf.Tg.SendEmailChanged(tg.Int64)
	}
}

func (ntf *CompositeNotifier) SendPasswordChanged(tx *AutoTx, userID *models.UserID) {
	const q = `
		SELECT email, verified, show_name, telegram
		FROM users
		WHERE id = $1
	`

	var email, name string
	var verified bool
	var tg sql.NullInt64
	tx.Query(q, userID.ID).Scan(&email, &verified, &name, &tg)

	if tx.Error() != nil {
		if tx.Error() != sql.ErrNoRows {
			ntf.srv.LogSystem().Error(tx.Error().Error())
		}
		return
	}

	if verified {
		ntf.Mail.SendPasswordChanged(email, name)
	}

	if tg.Valid {
		ntf.Tg.SendPasswordChanged(tg.Int64)
	}
}

func (ntf *CompositeNotifier) SendGreeting(address, showName string) {
	code := ntf.srv.TokenHash().VerificationCode(address)
	ntf.Mail.SendGreeting(address, showName, code)
}

func (ntf *CompositeNotifier) SendResetPassword(email, showName, gender string) {
	code, date := ntf.srv.TokenHash().ResetPasswordCode(email)
	ntf.Mail.SendResetPassword(email, showName, gender, code, date)
}

func (ntf *CompositeNotifier) entryTitle(tx *AutoTx, entryID int64) string {
	var title, content string
	tx.Query("SELECT title, edit_content FROM entries WHERE id = $1", entryID)
	tx.Scan(&title, &content)

	if title != "" {
		return title
	}

	content = ntf.md.RenderToString([]byte(content))
	content = RemoveHTML(content)
	content, _ = CutText(content, 100)

	return content
}

func (ntf *CompositeNotifier) SendNewComment(tx *AutoTx, cmt *models.Comment) {
	title := ntf.entryTitle(tx, cmt.EntryID)

	fromQ := sqlf.Select("gender.type, shadow_ban").
		From("users").
		Join("gender", "users.gender = gender.id").
		Where("users.id = ?", cmt.Author.ID)

	var fromGender string
	var fromShadowBan bool
	tx.QueryStmt(fromQ).Scan(&fromGender, &fromShadowBan)

	cmtIDQ := sqlf.Select("user_id").From("comments").Where("id = ?", cmt.ID)
	cmtUserID := tx.QueryStmt(cmtIDQ).ScanInt64()

	entryIDQ := sqlf.Select("user_id").From("entries").Where("id = ?", cmt.EntryID)
	entryUserID := tx.QueryStmt(entryIDQ).ScanInt64()

	toQ := sqlf.Select("users.id, users.name, show_name").
		Select("email, verified AND email_comments, telegram, telegram_comments").
		From("users").
		Join("watching", "watching.user_id = users.id").
		Where("watching.entry_id = ?", cmt.EntryID).
		Where("users.id <> ?", cmtUserID).
		Where("(? OR shadow_ban)", !fromShadowBan || cmtUserID == entryUserID)

	tx.QueryStmt(toQ)

	type userData struct {
		id        int64
		name      string
		showName  string
		email     string
		sendEmail bool
		sendTg    bool
		tg        sql.NullInt64
		canView   bool
	}

	var toUsers []*userData

	for {
		var user userData
		if !tx.Scan(&user.id, &user.name, &user.showName, &user.email, &user.sendEmail, &user.tg, &user.sendTg) {
			break
		}

		toUsers = append(toUsers, &user)
	}

	for _, user := range toUsers {
		userID, _ := LoadUserIDByID(tx, user.id)
		user.canView = CanViewEntry(tx, userID, cmt.EntryID)
		if !user.canView {
			continue
		}

		viewQ := sqlf.Select("TRUE").From("comments").Where("comments.id = ?", cmt.ID)
		AddViewCommentQuery(viewQ, userID)
		user.canView = tx.QueryStmt(viewQ).ScanBool()
	}

	for _, user := range toUsers {
		if !user.canView {
			continue
		}

		if user.sendEmail {
			ntf.Mail.SendNewComment(user.email, fromGender, user.showName, title, cmt)
		}

		if user.tg.Valid && user.sendTg {
			ntf.Tg.SendNewComment(user.tg.Int64, title, cmt)
		}

		ntf.Ntf.Notify(tx, cmt.ID, typeComment, user.name)
	}
}

func (ntf *CompositeNotifier) SendUpdateComment(tx *AutoTx, cmt *models.Comment) {
	title := ntf.entryTitle(tx, cmt.EntryID)

	ntf.Tg.SendUpdateComment(title, cmt)
	ntf.Ntf.NotifyUpdate(tx, cmt.ID, typeComment)
}

func (ntf *CompositeNotifier) SendRemoveComment(tx *AutoTx, commentID int64) {
	ntf.Tg.SendRemoveComment(commentID)
	ntf.Ntf.NotifyRemove(tx, commentID, typeComment)
}

func (ntf *CompositeNotifier) SendRead(name string, id int64) {
	ntf.Ntf.NotifyRead(name, id)
}

func (ntf *CompositeNotifier) SendInvited(tx *AutoTx, from, to string) {
	const toQ = `
		SELECT show_name, email, verified, telegram
		FROM users
		WHERE lower(name) = lower($1)
	`

	var sendEmail bool
	var toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, to).Scan(&toShowName, &email, &sendEmail, &tg)

	const fromQ = `
		SELECT users.id, show_name, gender.type
		FROM users, gender
		WHERE lower(users.name) = lower($1) AND users.gender = gender.id`

	var fromID int64
	var fromShowName, fromGender string
	tx.Query(fromQ, from).Scan(&fromID, &fromShowName, &fromGender)

	if sendEmail {
		ntf.Mail.SendInvited(email, fromShowName, fromGender, toShowName)
	}

	if tg.Valid {
		ntf.Tg.SendInvited(tg.Int64, from, fromShowName, fromGender)
	}

	ntf.Ntf.Notify(tx, fromID, typeInvited, to)
}

func (ntf *CompositeNotifier) SendNewFollower(tx *AutoTx, toPrivate bool, from, to string) {
	const toQ = `
		SELECT show_name, email, verified AND email_followers, telegram, telegram_followers, shadow_ban
		FROM users
		WHERE lower(name) = lower($1)
	`

	var sendEmail, sendTg, toShadowBan bool
	var toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, to).Scan(&toShowName, &email, &sendEmail, &tg, &sendTg, &toShadowBan)

	const fromQ = `
		SELECT users.id, show_name, gender.type, shadow_ban
		FROM users, gender
		WHERE lower(users.name) = lower($1) AND users.gender = gender.id`

	var fromID int64
	var fromShowName, fromGender string
	var fromShadowBan bool
	tx.Query(fromQ, from).Scan(&fromID, &fromShowName, &fromGender, &fromShadowBan)

	if fromShadowBan && !toShadowBan {
		return
	}

	if sendEmail {
		ntf.Mail.SendNewFollower(email, from, fromShowName, fromGender, toPrivate, toShowName)
	}

	if tg.Valid && sendTg {
		ntf.Tg.SendNewFollower(tg.Int64, to, from, fromShowName, fromGender, toPrivate)
	}

	ntf.Ntf.NotifyNewFollower(tx, fromID, to, toPrivate)
}

func (ntf *CompositeNotifier) SendRemoveFollower(tx *AutoTx, fromID int64, toName string) {
	const toQ = `
		SELECT id, telegram, telegram_followers
		FROM users
		WHERE lower(name) = lower($1)
	`

	var toID int64
	var sendTg bool
	var tg sql.NullInt64
	tx.Query(toQ, toName).Scan(&toID, &tg, &sendTg)

	const fromQ = `SELECT name FROM users WHERE id = $1`
	fromName := tx.QueryString(fromQ, fromID)

	if tg.Valid && sendTg {
		ntf.Tg.SendRemoveFollower(fromName, toName)
	}

	ntf.Ntf.NotifyRemoveFollower(tx, fromID, toID, toName)
}

func (ntf *CompositeNotifier) SendNewAccept(tx *AutoTx, from, to string) {
	const toQ = `
		SELECT show_name, email, verified AND email_followers, telegram, telegram_followers
		FROM users
		WHERE lower(name) = lower($1)
	`

	var sendEmail, sendTg bool
	var toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, to).Scan(&toShowName, &email, &sendEmail, &tg, &sendTg)

	const fromQ = `
		SELECT users.id, show_name, gender.type
		FROM users, gender
		WHERE lower(users.name) = lower($1) AND users.gender = gender.id`

	var fromID int64
	var fromShowName, fromGender string
	tx.Query(fromQ, from).Scan(&fromID, &fromShowName, &fromGender)

	if sendEmail {
		ntf.Mail.SendNewAccept(email, from, fromShowName, fromGender, toShowName)
	}

	if tg.Valid && sendTg {
		ntf.Tg.SendNewAccept(tg.Int64, from, fromShowName, fromGender)
	}

	ntf.Ntf.Notify(tx, fromID, typeAccept, to)
}

func (ntf *CompositeNotifier) SendNewCommentComplain(tx *AutoTx, commentID, fromID int64, content string) {
	const q = `
		SELECT entry_id, edit_content, author_id
		FROM comments
		WHERE comments.id = $1`

	var entryID, authorID int64
	var comment string
	tx.Query(q, commentID).Scan(&entryID, &comment, &authorID)

	from := LoadUser(tx, fromID)
	against := LoadUser(tx, authorID)

	ntf.Tg.SendCommentComplain(from, against, content, comment, commentID, entryID)
}

func (ntf *CompositeNotifier) SendNewEntryComplain(tx *AutoTx, entryID, fromID int64, content string) {
	const q = `
		SELECT edit_content, author_id
		FROM entries
		WHERE entries.id = $1`

	var entry string
	var authorID int64
	tx.Query(q, entryID).Scan(&entry, &authorID)

	from := LoadUser(tx, fromID)
	against := LoadUser(tx, authorID)

	ntf.Tg.SendEntryComplain(from, against, content, entry, entryID)
}

func (ntf *CompositeNotifier) SendNewMessageComplain(tx *AutoTx, msgID, fromID int64, content string) {
	const q = `
		SELECT edit_content, author_id
		FROM messages
		WHERE messages.id = $1`

	var msg string
	var authorID int64
	tx.Query(q, msgID).Scan(&msg, &authorID)

	from := LoadUser(tx, fromID)
	against := LoadUser(tx, authorID)

	ntf.Tg.SendMessageComplain(from, against, content, msg, msgID)
}

func (ntf *CompositeNotifier) SendNewUserComplain(tx *AutoTx, against *models.User, fromID int64, content string) {
	const q = `
		SELECT title
		FROM users
		WHERE users.id = $1`

	title := tx.QueryString(q, against.ID)
	from := LoadUser(tx, fromID)

	ntf.Tg.SendUserComplain(from, against, content, title)
}

func (ntf *CompositeNotifier) SendNewThemeComplain(tx *AutoTx, against *models.User, fromID int64, content string) {
	const q = `
		SELECT title
		FROM users
		WHERE users.id = $1`

	title := tx.QueryString(q, against.ID)
	from := LoadUser(tx, fromID)

	ntf.Tg.SendThemeComplain(from, against, content, title)
}

func (ntf *CompositeNotifier) SendNewWishComplain(tx *AutoTx, wishID, fromID int64) {
	const q = `
		SELECT content, from_id
		FROM wishes
		WHERE wishes.id = $1`

	var wish string
	var authorID int64
	tx.Query(q, wishID).Scan(&wish, &authorID)

	from := LoadUser(tx, fromID)
	against := LoadUser(tx, authorID)

	ntf.Tg.SendWishComplain(from, against, wish, wishID)
}

const retryQuery = `
SELECT EXISTS(SELECT 1
	FROM notifications
	JOIN notification_type ON notifications.type = notification_type.id
	JOIN users on user_id = users.id
	WHERE users.name = $1 AND notification_type.type = $2 AND age(notifications.created_at) < interval '6 month')
`

func (ntf *CompositeNotifier) SendAdmSent(tx *AutoTx, grandson string) {
	if tx.QueryBool(retryQuery, grandson, "adm_sent") {
		return
	}

	const toQ = `
		SELECT show_name, email, verified, telegram
		FROM users
		WHERE lower(name) = lower($1)
	`

	var sendEmail bool
	var toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, grandson).Scan(&toShowName, &email, &sendEmail, &tg)

	if sendEmail {
		ntf.Mail.SendAdmSent(email, toShowName)
	}

	if tg.Valid {
		ntf.Tg.SendAdmSent(tg.Int64)
	}

	ntf.Ntf.Notify(tx, 0, typeAdmSent, grandson)
}

func (ntf *CompositeNotifier) SendAdmReceived(tx *AutoTx, grandfather string) {
	if tx.QueryBool(retryQuery, grandfather, "adm_received") {
		return
	}

	const toQ = `
		SELECT show_name, email, verified, telegram
		FROM users
		WHERE lower(name) = lower($1)
	`

	var sendEmail bool
	var toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, grandfather).Scan(&toShowName, &email, &sendEmail, &tg)

	if sendEmail {
		ntf.Mail.SendAdmReceived(email, toShowName)
	}

	if tg.Valid {
		ntf.Tg.SendAdmReceived(tg.Int64)
	}

	ntf.Ntf.Notify(tx, 0, typeAdmReceived, grandfather)
}

func (ntf *CompositeNotifier) SendWishReceived(tx *AutoTx, wishID int64, user string) {
	ntf.Ntf.Notify(tx, wishID, typeWishReceived, user)
}

func (ntf *CompositeNotifier) SendWishCreated(tx *AutoTx, wishID int64, user string) {
	ntf.Ntf.Notify(tx, wishID, typeWishCreated, user)
}

func (ntf *CompositeNotifier) NotifyMessage(tx *AutoTx, msg *models.Message, user string) {
	const q = "SELECT telegram, telegram_messages FROM users WHERE lower(name) = lower($1)"

	var tg sql.NullInt64
	var sendTg bool
	tx.Query(q, user).Scan(&tg, &sendTg)

	if tg.Valid && sendTg {
		ntf.Tg.SendNewMessage(tg.Int64, msg)
	}

	ntf.Ntf.NotifyMessage(msg.ChatID, msg.ID, user)
}

func (ntf *CompositeNotifier) NotifyMessageUpdate(msg *models.Message, user string) {
	ntf.Tg.SendUpdateMessage(msg)
	ntf.Ntf.NotifyMessageUpdate(msg.ChatID, msg.ID, user)
}

func (ntf *CompositeNotifier) NotifyMessageRemove(chatID, msgID int64, user string) {
	ntf.Tg.SendRemoveMessage(msgID)
	ntf.Ntf.NotifyMessageRemove(chatID, msgID, user)
}

func (ntf *CompositeNotifier) NotifyMessageRead(chatID, msgID int64, user string) {
	ntf.Ntf.NotifyMessageRead(chatID, msgID, user)
}

func (ntf *CompositeNotifier) notifyMovedEntry(entryData []byte) {
	var entry models.Entry
	err := json.Unmarshal(entryData, &entry)
	if err != nil {
		ntf.srv.LogSystem().Error(err.Error())
		return
	}

	tx := NewAutoTx(ntf.srv.DB)
	defer tx.Finish()

	const toQ = `
		SELECT name, show_name,
			email, verified AND email_moved_entries,
			telegram, telegram_moved_entries
		FROM users
		WHERE id = $1
	`

	var sendEmail, sendTg bool
	var toName, toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, entry.User.ID).
		Scan(&toName, &toShowName, &email, &sendEmail, &tg, &sendTg)

	if sendEmail {
		ntf.Mail.SendEntryMoved(email, toShowName, entry.Title, entry.ID)
	}

	if tg.Valid && sendTg {
		ntf.Tg.SendEntryMoved(tg.Int64, entry.Title, entry.ID)
	}

	ntf.Ntf.Notify(tx, entry.ID, typeEntryMoved, toName)
}

func (ntf *CompositeNotifier) notifyBadge(badgeData []byte) {
	var badgeId struct {
		BadgeID int64 `json:"badge_id"`
		UserID  int64 `json:"user_id"`
	}
	err := json.Unmarshal(badgeData, &badgeId)
	if err != nil {
		ntf.srv.LogSystem().Error(err.Error())
		return
	}

	tx := NewAutoTx(ntf.srv.DB)
	defer tx.Finish()

	var badgeTitle, badgeDesc string
	const badgeQ = `
		SELECT title, description
		FROM badges
		WHERE id = $1
	`
	tx.Query(badgeQ, badgeId.BadgeID).Scan(&badgeTitle, &badgeDesc)

	const toQ = `
		SELECT name, show_name,
			email, verified AND email_badges,
			telegram, telegram_badges
		FROM users
		WHERE id = $1
	`

	var sendEmail, sendTg bool
	var toName, toShowName, email string
	var tg sql.NullInt64
	tx.Query(toQ, badgeId.UserID).
		Scan(&toName, &toShowName, &email, &sendEmail, &tg, &sendTg)

	if sendEmail {
		ntf.Mail.SendBadge(email, toName, toShowName, badgeTitle, badgeDesc)
	}

	if tg.Valid && sendTg {
		ntf.Tg.SendBadge(tg.Int64, toName, badgeTitle, badgeDesc)
	}

	ntf.Ntf.Notify(tx, badgeId.BadgeID, typeBadge, toName)
}
