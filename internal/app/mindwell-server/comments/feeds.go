package comments

import (
	"database/sql"
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/utils"
)

func baseFeedQuery(userID *models.UserID, limit int64) *sqlf.Stmt {
	return sqlf.Select("comments.id, entry_id").
		Select("extract(epoch from comments.created_at), comments.edit_content, comments.rating").
		Select("comments.up_votes, comments.down_votes, votes.vote").
		Select("comments.author_id, users.name, users.show_name").
		Select("is_online(users.last_seen_at), users.creator_id IS NOT NULL").
		Select("users.avatar, comments.user_id").
		From("comments").
		Join("users", "comments.author_id = users.id").
		With("votes",
			sqlf.Select("comment_id, vote").From("comment_votes").Where("comment_votes.user_id = ?", userID.ID)).
		LeftJoin("votes", "comments.id = votes.comment_id").
		Limit(limit)
}

func addEntryQuery(q *sqlf.Stmt, entryID int64) *sqlf.Stmt {
	return q.Where("entry_id = ?", entryID)
}

func entryFeedQuery(userID *models.UserID, entryID, limit int64) *sqlf.Stmt {
	q := baseFeedQuery(userID, limit)
	return addEntryQuery(q, entryID)
}

func addCanViewEntryQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	q.Join("entries", "entry_id = entries.id").
		Join("users AS authors", "entries.author_id = authors.id").
		Join("entry_privacy", "entries.visible_for = entry_privacy.id").
		Join("user_privacy", "authors.privacy = user_privacy.id")
	return utils.AddCanViewEntryQuery(q, userID)
}

func allFeedQuery(userID *models.UserID, limit int64) *sqlf.Stmt {
	q := baseFeedQuery(userID, limit)
	return addCanViewEntryQuery(q, userID)
}

func addUserQuery(q *sqlf.Stmt, userID *models.UserID, name string) *sqlf.Stmt {
	q.Where("lower(users.name) = lower(?)", name)
	return addCanViewEntryQuery(q, userID)
}

func userFeedQuery(userID *models.UserID, name string, limit int64) *sqlf.Stmt {
	q := baseFeedQuery(userID, limit)
	return addUserQuery(q, userID, name)
}

func addThemeQuery(q *sqlf.Stmt, userID *models.UserID, name string) *sqlf.Stmt {
	q.Where("lower(authors.name) = lower(?)", name)
	return addCanViewEntryQuery(q, userID)
}

func themeFeedQuery(userID *models.UserID, name string, limit int64) *sqlf.Stmt {
	q := baseFeedQuery(userID, limit)
	return addThemeQuery(q, userID, name)
}

func scrollQuery() *sqlf.Stmt {
	return sqlf.
		From("comments").
		Limit(1)
}

func addBordersQuery(q *sqlf.Stmt, before, after int64) *sqlf.Stmt {
	if after > 0 {
		return q.Where("comments.id > ?", after).
			OrderBy("comments.id ASC")
	} else if before > 0 {
		return q.Where("comments.id < ?", before).
			OrderBy("comments.id DESC")
	} else {
		return q.OrderBy("comments.id DESC")
	}
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

func loadNext(tx *utils.AutoTx, cmts *models.CommentList, scrollQ *sqlf.Stmt, before, after int64) {
	if len(cmts.Data) == 0 {
		scrollQ.Select("comments.id")

		if before > 0 {
			query := scrollQ.Clone().
				Where("comments.id > ?", before).
				OrderBy("comments.id ASC")
			nextAfter := tx.QueryStmt(query).ScanInt64()
			cmts.NextAfter = utils.FormatInt64(nextAfter)
		}

		if after > 0 {
			query := scrollQ.Clone().
				Where("comments.id < ?", after).
				OrderBy("comments.id DESC")
			nextBefore := tx.QueryStmt(query).ScanInt64()
			cmts.NextBefore = utils.FormatInt64(nextBefore)
		}
	} else {
		scrollQ.Select("TRUE")

		oldest := cmts.Data[0].ID
		newest := cmts.Data[len(cmts.Data)-1].ID

		if oldest > newest {
			newest, oldest = oldest, newest
		}

		beforeQ := scrollQ.Clone().Where("comments.id < ?", oldest)
		cmts.HasBefore = tx.QueryStmt(beforeQ).ScanBool()
		cmts.NextBefore = utils.FormatInt64(oldest)

		afterQ := scrollQ.Clone().Where("comments.id > ?", newest)
		cmts.HasAfter = tx.QueryStmt(afterQ).ScanBool()
		cmts.NextAfter = utils.FormatInt64(newest)
	}
}

// LoadEntryComments loads comments for entry.
func LoadEntryComments(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, entryID, limit int64, afterS, beforeS string) *models.CommentList {
	before := utils.ParseInt64(beforeS)
	after := utils.ParseInt64(afterS)

	query := entryFeedQuery(userID, entryID, limit)
	addBordersQuery(query, before, after)

	tx.QueryStmt(query)
	cmts := loadFeed(srv, tx, userID, after <= 0)

	scrollQ := scrollQuery()
	addEntryQuery(scrollQ, entryID)
	defer scrollQ.Close()

	loadNext(tx, cmts, scrollQ, before, after)

	return cmts
}

func loadUserComments(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, limit int64, author, afterS, beforeS string) *models.CommentList {
	before := utils.ParseInt64(beforeS)
	after := utils.ParseInt64(afterS)

	query := userFeedQuery(userID, author, limit)
	addBordersQuery(query, before, after)

	tx.QueryStmt(query)
	cmts := loadFeed(srv, tx, userID, after > 0)

	scrollQ := scrollQuery().
		Join("users", "comments.author_id = users.id")

	addUserQuery(scrollQ, userID, author)
	defer scrollQ.Close()

	loadNext(tx, cmts, scrollQ, before, after)

	return cmts
}

func loadThemeComments(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, limit int64, author, afterS, beforeS string) *models.CommentList {
	before := utils.ParseInt64(beforeS)
	after := utils.ParseInt64(afterS)

	query := themeFeedQuery(userID, author, limit)
	addBordersQuery(query, before, after)

	tx.QueryStmt(query)
	cmts := loadFeed(srv, tx, userID, after > 0)

	scrollQ := scrollQuery()
	addThemeQuery(scrollQ, userID, author)
	defer scrollQ.Close()

	loadNext(tx, cmts, scrollQ, before, after)

	return cmts
}

func loadAllComments(srv *utils.MindwellServer, tx *utils.AutoTx, userID *models.UserID, limit int64, afterS, beforeS string) *models.CommentList {
	before := utils.ParseInt64(beforeS)
	after := utils.ParseInt64(afterS)

	query := allFeedQuery(userID, limit)
	addBordersQuery(query, before, after)

	tx.QueryStmt(query)
	cmts := loadFeed(srv, tx, userID, after > 0)

	scrollQ := scrollQuery()
	addCanViewEntryQuery(scrollQ, userID)
	defer scrollQ.Close()

	loadNext(tx, cmts, scrollQ, before, after)

	return cmts
}
