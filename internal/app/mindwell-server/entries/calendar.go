package entries

import (
	"database/sql"
	"github.com/go-openapi/runtime/middleware"
	"github.com/leporo/sqlf"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/entries"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/themes"
	"github.com/sevings/mindwell-server/restapi/operations/users"
	"github.com/sevings/mindwell-server/utils"
	"strings"
	"time"
)

func baseCalendarQuery(cal *models.Calendar) *sqlf.Stmt {
	q := sqlf.Select("entries.id, extract(epoch from entries.created_at) as created_at").
		Select("entries.title, entries.edit_content").
		From("entries")

	if cal.Start > 0 {
		q.Where("entries.created_at >= to_timestamp(?)", cal.Start)
	}

	if cal.End > 0 {
		q.Where("entries.created_at < to_timestamp(?)", cal.End)
	}

	if cal.Start > 0 {
		q.OrderBy("entries.created_at ASC")
	} else {
		q.OrderBy("entries.created_at DESC")
	}

	return q.Limit(cal.Limit)
}

func myCalendarQuery(userID *models.UserID, cal *models.Calendar) *sqlf.Stmt {
	return baseCalendarQuery(cal).
		Where("author_id = ?", userID.ID)
}

func tlogCalendarQuery(userID *models.UserID, tlog string, cal *models.Calendar) *sqlf.Stmt {
	q := baseCalendarQuery(cal).
		Join("users AS authors", "entries.author_id = authors.id").
		Join("users AS entry_users", "entries.user_id = entry_users.id").
		Where("lower(authors.name) = lower(?)", tlog).
		Where(`(entries.user_id = entries.author_id OR NOT entry_users.shadow_ban OR ?)`, userID.Ban.Shadow).
		Join("entry_privacy", "entries.visible_for = entry_privacy.id")
	AddRelationToTlogQuery(q, userID, tlog)
	return utils.AddEntryOpenQuery(q, userID, false)
}

func loadEmptyCalendar(tx *utils.AutoTx, q *sqlf.Stmt, start, end, limit int64) *models.Calendar {
	var createdAt float64
	tx.QueryStmt(q)
	tx.Scan(&createdAt)

	const maxDuration = 60*60*24*7*6 + 1 // six weeks
	minDate := int64(createdAt)
	maxDate := time.Now().Unix() + 1

	if start > 0 && start < minDate {
		start = minDate
	}
	if end > 0 && end > maxDate {
		end = maxDate
	}

	if start > 0 && end > 0 && end-start > maxDuration {
		end = start + maxDuration
	}

	return &models.Calendar{
		Start: start,
		End:   end,
		Limit: limit,
	}
}

func loadCalendarEntry(tx *utils.AutoTx) *models.CalendarEntry {
	var title, content string
	var entry models.CalendarEntry
	ok := tx.Scan(&entry.ID, &entry.CreatedAt,
		&title, &content)
	if !ok {
		return nil
	}

	title = strings.TrimSpace(title)
	if title != "" {
		title = bluemonday.StrictPolicy().Sanitize(title)
		entry.Title, _ = utils.CutText(title, 100)
	} else {
		content = strings.TrimSpace(content)
		content = md.RenderToString([]byte(content))
		content = utils.RemoveHTML(content)
		entry.Title, _ = utils.CutHtml(content, 1, 100, 0)
	}

	return &entry
}

func loadCalendar(tx *utils.AutoTx, cal *models.Calendar) {
	for {
		entry := loadCalendarEntry(tx)
		if entry == nil {
			break
		}

		cal.Entries = append(cal.Entries, entry)
	}

	if cal.Start > 0 {
		list := cal.Entries
		for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
			list[i], list[j] = list[j], list[i]
		}
	}
}

func loadTlogCalendar(tx *utils.AutoTx, userID *models.UserID, tlog string, start, end, limit int64) *models.Calendar {
	if userID.Name == tlog {
		return loadMyCalendar(tx, userID, start, end, limit)
	}

	createdAtQuery := sqlf.Select("extract(epoch FROM created_at)").
		From("users").
		Where("lower(name) = lower(?)", tlog)

	cal := loadEmptyCalendar(tx, createdAtQuery, start, end, limit)
	if cal.End > 0 && cal.Start >= cal.End {
		return cal
	}

	q := tlogCalendarQuery(userID, tlog, cal)
	tx.QueryStmt(q)
	loadCalendar(tx, cal)

	return cal
}

func newTlogCalendarLoader(srv *utils.MindwellServer) func(users.GetUsersNameCalendarParams, *models.UserID) middleware.Responder {
	return func(params users.GetUsersNameCalendarParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_tlog")
				return users.NewGetUsersNameCalendarNotFound().WithPayload(err)
			}

			feed := loadTlogCalendar(tx, userID, params.Name, *params.Start, *params.End, *params.Limit)
			return users.NewGetUsersNameCalendarOK().WithPayload(feed)
		})
	}
}

func newThemeCalendarLoader(srv *utils.MindwellServer) func(themes.GetThemesNameCalendarParams, *models.UserID) middleware.Responder {
	return func(params themes.GetThemesNameCalendarParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewTlogName(tx, userID, params.Name)
			if !canView {
				err := srv.StandardError("no_theme")
				return themes.NewGetThemesNameCalendarNotFound().WithPayload(err)
			}

			feed := loadTlogCalendar(tx, userID, params.Name, *params.Start, *params.End, *params.Limit)
			return themes.NewGetThemesNameCalendarOK().WithPayload(feed)
		})
	}
}

func loadMyCalendar(tx *utils.AutoTx, userID *models.UserID, start, end, limit int64) *models.Calendar {
	createdAtQuery := sqlf.Select("extract(epoch FROM created_at)").
		From("users").
		Where("id = ?", userID.ID)

	cal := loadEmptyCalendar(tx, createdAtQuery, start, end, limit)
	if cal.End > 0 && cal.Start >= cal.End {
		return cal
	}

	q := myCalendarQuery(userID, cal)
	tx.QueryStmt(q)
	loadCalendar(tx, cal)

	return cal
}

func newMyCalendarLoader(srv *utils.MindwellServer) func(me.GetMeCalendarParams, *models.UserID) middleware.Responder {
	return func(params me.GetMeCalendarParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			feed := loadMyCalendar(tx, userID, *params.Start, *params.End, *params.Limit)

			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(nil)
				return me.NewPutMeCoverBadRequest().WithPayload(err)
			}

			return me.NewGetMeCalendarOK().WithPayload(feed)
		})
	}
}

func loadAdjacent(tx *utils.AutoTx, userID *models.UserID, entryID int64) models.AdjacentEntries {
	var createdAt time.Time
	var authorID int64
	curQuery := sqlf.Select("author_id, created_at").From("entries").Where("id = ?", entryID)
	tx.QueryStmt(curQuery).Scan(&authorID, &createdAt)

	q := sqlf.Select("entries.id, extract(epoch from entries.created_at) as created_at").
		Select("entries.title, entries.edit_content").
		From("entries").
		Where("author_id = ?", authorID).
		Limit(1)

	if authorID != userID.ID {
		q.Join("entry_privacy", "entries.visible_for = entry_privacy.id")
		q.Join("users AS authors", "entries.author_id = authors.id")
		q.Join("users AS entry_users", "entries.user_id = entry_users.id")
		q.Where(`(entries.user_id = entries.author_id OR NOT entry_users.shadow_ban OR ?)`, userID.Ban.Shadow)
		q = addRelationToQuery(q, userID, authorID)
		q = utils.AddEntryOpenQuery(q, userID, false)
	}

	newerQuery := q.Clone().
		Where("entries.created_at > ?", createdAt).
		OrderBy("entries.created_at ASC")
	tx.QueryStmt(newerQuery)
	newer := loadCalendarEntry(tx)

	olderQuery := q.Clone().
		Where("entries.created_at < ?", createdAt).
		OrderBy("entries.created_at DESC")
	tx.QueryStmt(olderQuery)
	older := loadCalendarEntry(tx)

	return models.AdjacentEntries{
		ID:    entryID,
		Newer: newer,
		Older: older,
	}
}

func newAdjacentLoader(srv *utils.MindwellServer) func(entries.GetEntriesIDAdjacentParams, *models.UserID) middleware.Responder {
	return func(params entries.GetEntriesIDAdjacentParams, userID *models.UserID) middleware.Responder {
		return srv.Transact(func(tx *utils.AutoTx) middleware.Responder {
			canView := utils.CanViewEntry(tx, userID, params.ID)
			if !canView {
				err := srv.StandardError("no_entry")
				return entries.NewGetEntriesIDAdjacentNotFound().WithPayload(err)
			}

			adj := loadAdjacent(tx, userID, params.ID)
			if tx.Error() != nil && tx.Error() != sql.ErrNoRows {
				err := srv.NewError(nil)
				return me.NewPutMeCoverBadRequest().WithPayload(err)
			}

			return entries.NewGetEntriesIDAdjacentOK().WithPayload(&adj)
		})
	}
}
