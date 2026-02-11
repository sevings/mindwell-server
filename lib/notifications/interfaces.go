package notifications

import "github.com/sevings/mindwell-server/models"

// MailSender interface for email notification services
type MailSender interface {
	SendGreeting(address, name, code string)
	SendPasswordChanged(address, name string)
	SendEmailChanged(address, name string)
	SendResetPassword(address, name, gender, code string, date int64)
	SendNewComment(address, fromGender, toShowName, entryTitle string, cmt *models.Comment)
	SendNewFollower(address, fromName, fromShowName, fromGender string, toPrivate bool, toShowName string)
	SendNewAccept(address, fromName, fromShowName, fromGender, toShowName string)
	SendNewInvite(address, name string)
	SendInvited(address, fromShowName, fromGender, toShowName string)
	SendAdmSent(address, toShowName string)
	SendAdmReceived(address, toShowName string)
	SendCommentComplain(from, against, content, comment string, commentID, entryID int64)
	SendEntryComplain(from, against, content, entry string, entryID int64)
	SendEntryMoved(address, toShowName, entryTitle string, entryID int64)
	SendBadge(address, toName, toShowName, badgeTitle, badgeDesc string)
	Stop()
}
