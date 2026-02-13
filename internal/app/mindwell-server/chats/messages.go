package chats

import (
	"time"

	"github.com/sevings/mindwell-server/lib/database"
	"github.com/sevings/mindwell-server/lib/server"

	"github.com/go-openapi/runtime/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/comments"
	"github.com/sevings/mindwell-server/internal/app/mindwell-server/users"
	"github.com/sevings/mindwell-server/lib/helpers"
	"github.com/sevings/mindwell-server/lib/userutil"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/chats"
)

const loadMessagesQuery = `
    SELECT messages.id, extract(epoch from created_at), author_id,
        content, edit_content
    FROM messages
    WHERE chat_id = $1
`

func loadMessageList(srv *server.MindwellServer, tx *database.AutoTx, userID *models.UserID, chatID int64, reverse bool) *models.MessageList {
	var result models.MessageList

	for {
		msg := models.Message{
			Author: &models.User{},
			ChatID: chatID,
		}
		ok := tx.Scan(&msg.ID, &msg.CreatedAt, &msg.Author.ID,
			&msg.Content, &msg.EditContent)

		if !ok {
			break
		}

		if msg.Author.ID == userID.ID {
			msg.Rights = &models.MessageRights{
				Delete: true,
				Edit:   true,
			}
		} else {
			msg.EditContent = ""
			msg.Rights = &models.MessageRights{
				Complain: !userID.Ban.Complain,
			}
		}

		result.Data = append(result.Data, &msg)
	}

	talkers := make(map[int64]*models.User, 2)

	for _, msg := range result.Data {
		author, found := talkers[msg.Author.ID]
		if !found {
			author = users.LoadUserByID(srv, tx, msg.Author.ID)
			talkers[msg.Author.ID] = author
		}
		msg.Author = author
	}

	if reverse {
		list := result.Data
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
	}

	return &result
}

const userLastReadQuery = "SELECT last_read FROM talkers WHERE chat_id = $1 AND user_id = $2"
const partnerLastReadQuery = `
	SELECT COALESCE(
		(SELECT MAX(last_read) FROM talkers WHERE chat_id = $1 AND user_id <> $2),
		12147483647
	)`

func setMessagesRead(tx *database.AutoTx, list *models.MessageList, userID int64) {
	chatID := list.Data[0].ChatID
	userLastRead := tx.QueryInt64(userLastReadQuery, chatID, userID)
	partnerLastRead := tx.QueryInt64(partnerLastReadQuery, chatID, userID)
	for _, msg := range list.Data {
		if msg.Author.ID == userID {
			msg.Read = msg.ID <= partnerLastRead
		} else {
			msg.Read = msg.ID <= userLastRead
		}
	}
}

func newMessageListLoader(srv *server.MindwellServer) func(chats.GetChatsNameMessagesParams, *models.UserID) middleware.Responder {
	return func(params chats.GetChatsNameMessagesParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
			chatID, partnerID := findDialog(tx, userID.ID, params.Name)
			if partnerID == 0 {
				err := srv.StandardError("no_chat")
				return chats.NewGetChatsNameMessagesNotFound().WithPayload(err)
			}
			if chatID == 0 {
				return chats.NewGetChatsNameMessagesOK()
			}

			var q = loadMessagesQuery

			before := helpers.ParseInt64(*params.Before)
			after := helpers.ParseInt64(*params.After)

			if after > 0 {
				q := q + " AND messages.id > $2 ORDER BY messages.id ASC LIMIT $3"
				tx.Query(q, chatID, after, *params.Limit)
			} else if before > 0 {
				q := q + " AND messages.id < $2 ORDER BY messages.id DESC LIMIT $3"
				tx.Query(q, chatID, before, *params.Limit)
			} else {
				q := q + " ORDER BY messages.id DESC LIMIT $2"
				tx.Query(q, chatID, *params.Limit)
			}

			list := loadMessageList(srv, tx, userID, chatID, after <= 0)

			if len(list.Data) == 0 {
				return chats.NewGetChatsNameMessagesOK().WithPayload(list)
			}

			setMessagesRead(tx, list, userID.ID)

			const unreadQuery = "SELECT unread_count FROM talkers WHERE chat_id = $1 AND user_id = $2"
			list.UnreadCount = tx.QueryInt64(unreadQuery, chatID, userID.ID)

			const beforeQuery = `SELECT EXISTS(
				SELECT 1
				FROM messages
                WHERE chat_id = $1 AND id < $2)`

			nextBefore := list.Data[0].ID
			list.NextBefore = helpers.FormatInt64(nextBefore)
			tx.Query(beforeQuery, chatID, nextBefore)
			tx.Scan(&list.HasBefore)

			const afterQuery = `SELECT EXISTS(
				SELECT 1
				FROM messages
                WHERE chat_id = $1 AND id > $2)`

			nextAfter := list.Data[len(list.Data)-1].ID
			list.NextAfter = helpers.FormatInt64(nextAfter)
			tx.Query(afterQuery, chatID, nextAfter)
			tx.Scan(&list.HasAfter)

			return chats.NewGetChatsNameMessagesOK().WithPayload(list)
		})
	}
}

func canSendMessage(tx *database.AutoTx, userID *models.UserID, partnerID, chatID int64) bool {
	if userID.ID == partnerID {
		return true
	}

	toMe := userutil.LoadRelation(tx, partnerID, userID.ID)
	if toMe == models.RelationshipRelationIgnored {
		return false
	}

	const lastMessageQuery = `SELECT last_message FROM chats WHERE id = $1`
	lastMessage := tx.QueryInt64(lastMessageQuery, chatID)
	if lastMessage != 0 {
		return true
	}

	const chatPrivacyQuery = `
SELECT type, shadow_ban
FROM users
JOIN user_chat_privacy ON users.chat_privacy = user_chat_privacy.id
WHERE users.id = $1`
	var chatPrivacy string
	var partnerShadowBan bool
	tx.Query(chatPrivacyQuery, partnerID).Scan(&chatPrivacy, &partnerShadowBan)

	switch chatPrivacy {
	case "invited":
		return userID.IsInvited && (!userID.Ban.Shadow || partnerShadowBan)
	case "followers":
		return userID.IsInvited && (!userID.Ban.Shadow || partnerShadowBan) &&
			userutil.LoadRelation(tx, userID.ID, partnerID) == models.RelationshipRelationFollowed
	case "friends":
		return toMe == models.RelationshipRelationFollowed &&
			userutil.LoadRelation(tx, userID.ID, partnerID) == models.RelationshipRelationFollowed
	case "me":
		return false
	}

	return false
}

const createMessageQuery = `
    INSERT INTO messages(chat_id, author_id, content, edit_content)
    VALUES($1, $2, $3, $4)
    RETURNING id, extract(epoch from created_at)
`

func createMessage(srv *server.MindwellServer, tx *database.AutoTx, userID, chatID int64, content string) *models.Message {
	msg := &models.Message{
		ChatID:      chatID,
		Author:      users.LoadUserByID(srv, tx, userID),
		Content:     comments.HtmlContent(content),
		EditContent: content,
		Rights: &models.MessageRights{
			Delete: true,
			Edit:   true,
		},
	}

	tx.Query(createMessageQuery, chatID, userID, msg.Content, msg.EditContent)
	tx.Scan(&msg.ID, &msg.CreatedAt)

	setMessageRead(tx, msg, userID)

	return msg
}

func newMessageCreator(srv *server.MindwellServer) func(chats.PostChatsNameMessagesParams, *models.UserID) middleware.Responder {
	cantChatErr := srv.NewError(&i18n.Message{ID: "cant_chat", Other: "You are not allowed to send messages to this chat."})
	return func(params chats.PostChatsNameMessagesParams, userID *models.UserID) middleware.Responder {
		if userID.ID == 0 {
			return chats.NewPostChatsNameMessagesForbidden().WithPayload(cantChatErr)
		}

		return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
			msg := getCachedMessage(userID.ID, params.UID, params.Name)
			if msg != nil {
				return chats.NewPostChatsNameMessagesCreated().WithPayload(msg)
			}

			chatID, partnerID := findDialog(tx, userID.ID, params.Name)
			if partnerID == 0 {
				err := srv.StandardError("no_chat")
				return chats.NewPostChatsNameMessagesNotFound().WithPayload(err)
			}

			if chatID == 0 {
				chatID = createChat(srv, tx, userID.ID, partnerID).ID
			}

			if !canSendMessage(tx, userID, partnerID, chatID) {
				return chats.NewPostChatsNameMessagesForbidden().WithPayload(cantChatErr)
			}

			msg = createMessage(srv, tx, userID.ID, chatID, params.Content)

			const q = "UPDATE talkers SET last_read = $3, unread_count = 0 WHERE chat_id = $1 AND user_id = $2"
			tx.Exec(q, chatID, userID.ID, msg.ID)

			setCachedMessage(userID.ID, params.UID, params.Name, msg)
			return chats.NewPostChatsNameMessagesCreated().WithPayload(msg)
		})
	}
}

const loadMessageQuery = `
    SELECT chat_id, extract(epoch from messages.created_at),
        content, edit_content,
        users.id, users.name, users.show_name,
        is_online(users.last_seen_at), users.avatar
    FROM messages
    JOIN users ON users.id = messages.author_id
    WHERE messages.id = $1
`

func loadMessage(srv *server.MindwellServer, tx *database.AutoTx, userID *models.UserID, msgID int64) *models.Message {
	var avatar string
	msg := &models.Message{
		ID:     msgID,
		Author: &models.User{},
	}

	tx.Query(loadMessageQuery, msgID)
	tx.Scan(&msg.ChatID, &msg.CreatedAt,
		&msg.Content, &msg.EditContent,
		&msg.Author.ID, &msg.Author.Name, &msg.Author.ShowName,
		&msg.Author.IsOnline, &avatar)

	msg.Author.Avatar = srv.NewAvatar(avatar)

	if msg.Author.ID == userID.ID {
		msg.Rights = &models.MessageRights{
			Delete: true,
			Edit:   true,
		}
	} else {
		msg.EditContent = ""
		msg.Rights = &models.MessageRights{
			Complain: !userID.Ban.Complain,
		}
	}

	return msg
}

func canViewChat(tx *database.AutoTx, userID, chatID int64) bool {
	const q = "SELECT true FROM talkers WHERE user_id = $1 AND chat_id = $2"
	return tx.QueryBool(q, userID, chatID)
}

func setMessageRead(tx *database.AutoTx, msg *models.Message, userID int64) {
	var lastRead int64
	if msg.Author.ID == userID {
		lastRead = tx.QueryInt64(partnerLastReadQuery, msg.ChatID, userID)
	} else {
		lastRead = tx.QueryInt64(userLastReadQuery, msg.ChatID, userID)
	}

	msg.Read = msg.ID <= lastRead
}

func newMessageLoader(srv *server.MindwellServer) func(chats.GetMessagesIDParams, *models.UserID) middleware.Responder {
	return func(params chats.GetMessagesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
			msg := loadMessage(srv, tx, userID, params.ID)
			if msg.CreatedAt == 0 {
				err := srv.StandardError("no_message")
				return chats.NewGetMessagesIDNotFound().WithPayload(err)
			}

			if !canViewChat(tx, userID.ID, msg.ChatID) {
				err := srv.StandardError("no_message")
				return chats.NewGetMessagesIDForbidden().WithPayload(err)
			}

			setMessageRead(tx, msg, userID.ID)

			return chats.NewGetMessagesIDOK().WithPayload(msg)
		})
	}
}

func canEditMessage(tx *database.AutoTx, userID, msgID int64) bool {
	const q = "SELECT author_id = $1 FROM messages WHERE id = $2"
	return tx.QueryBool(q, userID, msgID)
}

const editMessageQuery = `
    UPDATE messages
    SET content = $2, edit_content = $3
    WHERE id = $1
    RETURNING chat_id, extract(epoch from created_at)
`

func editMessage(srv *server.MindwellServer, tx *database.AutoTx, userID, msgID int64, content string) *models.Message {
	msg := &models.Message{
		ID:          msgID,
		Author:      users.LoadUserByID(srv, tx, userID),
		Content:     comments.HtmlContent(content),
		EditContent: content,
		Rights: &models.MessageRights{
			Delete: true,
			Edit:   true,
		},
	}

	tx.Query(editMessageQuery, msgID, msg.Content, msg.EditContent)
	tx.Scan(&msg.ChatID, &msg.CreatedAt)

	return msg
}

func newMessageEditor(srv *server.MindwellServer) func(chats.PutMessagesIDParams, *models.UserID) middleware.Responder {
	return func(params chats.PutMessagesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
			if !canEditMessage(tx, userID.ID, params.ID) {
				err := srv.StandardError("no_message")
				return chats.NewPutMessagesIDForbidden().WithPayload(err)
			}

			msg := editMessage(srv, tx, userID.ID, params.ID, params.Content)
			if msg.CreatedAt == 0 {
				err := srv.StandardError("no_message")
				return chats.NewPutMessagesIDNotFound().WithPayload(err)
			}

			setMessageRead(tx, msg, userID.ID)

			return chats.NewPutMessagesIDOK().WithPayload(msg)
		})
	}
}

func deleteMessage(tx *database.AutoTx, msgID int64) int64 {
	const q = "DELETE FROM messages WHERE id = $1 RETURNING chat_id"
	return tx.QueryInt64(q, msgID)
}

func newMessageDeleter(srv *server.MindwellServer) func(chats.DeleteMessagesIDParams, *models.UserID) middleware.Responder {
	return func(params chats.DeleteMessagesIDParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *database.AutoTx) middleware.Responder {
			if !canEditMessage(tx, userID.ID, params.ID) {
				err := srv.StandardError("no_message")
				return chats.NewDeleteMessagesIDForbidden().WithPayload(err)
			}

			deleteMessage(tx, params.ID)

			return chats.NewDeleteMessagesIDOK()
		})
	}
}

func SendWelcomeMessage(srv *server.MindwellServer, user *models.AuthProfile) {
	helpURL := srv.ConfigString("server.base_url") + "help/faq"

	text := `Привет, друг! Мы рады видеть тебя с нами!
У нас уютно. Убедись в этом лично, написав первый пост в своем дневнике.
На данный момент тебе доступны основные функции сайта. Продолжай публиковать открытые посты, чтобы получить приглашение и иметь возможность комментировать записи других пользователей, голосовать, начинать новые беседы и многое другое. Также ты можешь обратиться за приглашением ко мне.
Ответы на распространенные вопросы содержатся в разделе Помощь ` + helpURL + `
Если в разделе ответа не нашлось, спрашивай у меня.
Чувствуй себя как дома 😌`

	tx := database.NewAutoTx(srv.DB)
	defer tx.Finish()

	chat := createChat(srv, tx, 1, user.ID)
	msg := createMessage(srv, tx, 1, chat.ID, text)

	setCachedMessage(1, time.Now().Unix(), user.Name, msg)
}
