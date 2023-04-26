package comments

import (
	"database/sql"
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
)

func baseFeedQuery(userID *models.UserID) *sqlf.Stmt {
	return sqlf.Select("comments.id, entry_id").
		Select("extract(epoch from comments.created_at), edit_content, rating").
		Select("up_votes, down_votes, votes.vote").
		Select("author_id, name, show_name").
		Select("is_online(last_seen_at), users.creator_id IS NOT NULL").
		Select("avatar, user_id").
		From("comments").
		Join("users", "comments.author_id = users.id").
		With("votes",
			sqlf.Select("comment_id, vote").From("comment_votes").Where("user_id = ?", userID.ID)).
		LeftJoin("votes", "comments.id = votes.comment_id ")
}

func loadFeed(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, reverse bool) *models.CommentList {
	var list []*models.Comment
	var cmtUserIDs []int64

	for {
		comment := models.Comment{
			Author: &models.User{},
			Rating: &models.Rating{
				IsVotable: true,
			},
		}
		var vote sql.NullFloat64
		var avatar string
		var cmtUserID int64
		ok := tx.Scan(&comment.ID, &comment.EntryID,
			&comment.CreatedAt, &comment.EditContent, &comment.Rating.Rating,
			&comment.Rating.UpCount, &comment.Rating.DownCount, &vote,
			&comment.Author.ID, &comment.Author.Name, &comment.Author.ShowName,
			&comment.Author.IsOnline, &comment.Author.IsTheme,
			&avatar, &cmtUserID)
		if !ok {
			break
		}

		cmtUserIDs = append(cmtUserIDs, cmtUserID)

		setCommentText(&comment)

		if comment.Author.IsTheme {
			comment.Author.IsOnline = false
		}

		comment.Rating.Vote = commentVote(vote)
		comment.Rating.ID = comment.ID
		comment.Author.Avatar = srv.NewAvatar(avatar)
		list = append(list, &comment)
	}

	entryUsers := make(map[int64]int64)
	themeCreators := make(map[int64]int64)

	for i := 0; i < len(list); i++ {
		cmt := list[i]
		cmtUserID := cmtUserIDs[i]
		entryID := cmt.EntryID

		var themeCreatorID int64
		entryUserID, ok := entryUsers[entryID]
		if ok {
			themeCreatorID = themeCreators[entryID]
		} else {
			entryUserID, themeCreatorID = LoadEntryAuthor(tx, entryID)
			entryUsers[entryID] = entryUserID
			themeCreators[entryID] = themeCreatorID
		}

		setCommentRights(cmt, userID, cmtUserID, entryUserID, themeCreatorID)

		if !cmt.Rights.Edit {
			cmt.EditContent = ""
		}
	}

	if reverse {
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
	}

	return &models.CommentList{Data: list}
}
