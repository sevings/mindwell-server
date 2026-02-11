package userutil

import (
	"github.com/leporo/sqlf"
	"github.com/sevings/mindwell-server/models"
)

// AddAuthorRelationsToMeQuery adds a CTE for relations from authors to the current user.
func AddAuthorRelationsToMeQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	return q.
		With("relations_to_me",
			sqlf.Select("relation.type, relations.from_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.to_id = ?", userID.ID)).
		LeftJoin("relations_to_me", "relations_to_me.from_id = authors.id")
}

// AddAuthorRelationsFromMeQuery adds a CTE for relations from the current user to authors.
func AddAuthorRelationsFromMeQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	return q.
		With("relations_from_me",
			sqlf.Select("relation.type, relations.to_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.from_id = ?", userID.ID)).
		LeftJoin("relations_from_me", "relations_from_me.to_id = authors.id")
}

// AddViewCommentQuery adds privacy filtering for viewing comments.
func AddViewCommentQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	return q.
		Join("entries", "entry_id = entries.id").
		With("relations_to_me",
			sqlf.Select("relation.type, relations.from_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.to_id = ?", userID.ID)).
		LeftJoin("relations_to_me", "relations_to_me.from_id = comments.author_id").
		Where("(relations_to_me.type IS NULL OR relations_to_me.type <> 'ignored' OR entries.user_id = ?)", userID.ID).
		With("relations_from_me",
			sqlf.Select("relation.type, relations.to_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.from_id = ?", userID.ID)).
		LeftJoin("relations_from_me", "relations_from_me.to_id = comments.author_id").
		Where("(relations_from_me.type IS NULL OR relations_from_me.type <> 'ignored' OR comments.author_id = entries.author_id)")
}

// AddEntryOpenQuery adds privacy filtering for entry visibility based on entry privacy settings.
func AddEntryOpenQuery(q *sqlf.Stmt, userID *models.UserID, showMe bool) *sqlf.Stmt {
	return q.Where(`
CASE entry_privacy.type
WHEN 'all' THEN TRUE
WHEN 'registered' THEN ?
WHEN 'invited' THEN (? AND (NOT ? OR entry_users.shadow_ban)) OR authors.id = ?
WHEN 'followers' THEN authors.id = ? OR authors.creator_id = ? OR relations_from_me.type = 'followed'
WHEN 'some' THEN authors.id = ? OR authors.creator_id = ?
	OR EXISTS(SELECT 1 from entries_privacy WHERE user_id = ? AND entry_id = entries.id)
WHEN 'me' THEN ? AND authors.id = ?
ELSE FALSE
END
`, userID.ID > 0, userID.IsInvited, userID.Ban.Shadow, userID.ID, userID.ID, userID.ID, userID.ID, userID.ID, userID.ID, showMe, userID.ID)
}

// AddCanViewAuthorQuery adds privacy filtering for checking if a user can view an author's profile.
func AddCanViewAuthorQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddAuthorRelationsFromMeQuery(q, userID)
	AddAuthorRelationsToMeQuery(q, userID)

	return q.
		Where("(relations_to_me.type IS NULL OR relations_to_me.type <> 'ignored')").
		Where(`(
authors.id = ? OR
CASE user_privacy.type
WHEN 'all' THEN TRUE
WHEN 'registered' THEN ?
WHEN 'invited' THEN ? AND (NOT ? OR authors.shadow_ban)
WHEN 'followers' THEN relations_from_me.type = 'followed'
ELSE FALSE
END)`, userID.ID, userID.ID > 0, userID.IsInvited, userID.Ban.Shadow)
}

// AddViewAuthorQuery adds privacy filtering for viewing an author (includes ignore check).
func AddViewAuthorQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddCanViewAuthorQuery(q, userID)
	return q.Where("(relations_from_me.type IS NULL OR relations_from_me.type <> 'ignored')")
}

// AddViewEntryQuery adds full privacy filtering for viewing entries.
func AddViewEntryQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddEntryOpenQuery(q, userID, true)
	return AddViewAuthorQuery(q, userID)
}

// AddCanViewEntryQuery adds privacy filtering for checking if a user can view an entry.
func AddCanViewEntryQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddEntryOpenQuery(q, userID, true)
	return AddCanViewAuthorQuery(q, userID)
}
