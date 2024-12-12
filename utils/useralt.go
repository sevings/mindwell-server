package utils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/leporo/sqlf"
	"github.com/mileusna/useragent"
)

type UserAltQuery struct {
	tx      *AutoTx
	limit   int
	baseURL string
	ipAPI   string
}

type UserAltCount struct {
	Count   int64
	Alt     string
	baseURL string
}

type UserIPCount struct {
	Count int64
	From  time.Time
	To    time.Time
	IP    string
	ipAPI string
}

type UserAppCount struct {
	Count     int64
	From      time.Time
	To        time.Time
	UserAgent string
}

type UserAltCounts []UserAltCount
type UserIPCounts []UserIPCount
type UserAppCounts []UserAppCount

func NewUserAltQuery(tx *AutoTx, limit int, baseURL, ipAPI string) *UserAltQuery {
	return &UserAltQuery{
		tx:      tx,
		limit:   limit,
		baseURL: baseURL,
		ipAPI:   ipAPI,
	}
}

func (q *UserAltQuery) GetSuspectedAlt(user string) (string, bool) {
	var alt sql.NullString
	var conf bool
	confQuery := sqlf.Select("alt_of, confirmed_alt").
		From("users").
		Where("lower(name) = lower(?)", user)
	q.tx.QueryStmt(confQuery).Scan(&alt, &conf)
	return alt.String, conf
}

func (q *UserAltQuery) GetIPAlts(user string) UserAltCounts {
	return q.getAlts(q.buildIPQuery(user))
}

func (q *UserAltQuery) GetDeviceAlts(user string) UserAltCounts {
	return q.getAlts(q.buildDeviceQuery(user))
}

func (q *UserAltQuery) GetAppAlts(user string) UserAltCounts {
	return q.getAlts(q.buildAppQuery(user))
}

func (q *UserAltQuery) GetUID2Alts(user string) UserAltCounts {
	return q.getAlts(q.buildUID2Query(user))
}

func (q *UserAltQuery) GetEmailAlts(user, email string) UserAltCounts {
	return q.getAlts(q.buildEmailQuery(user, email))
}

func (q *UserAltQuery) GetCommonIPs(userA, userB string) UserIPCounts {
	return q.getIPCounts(q.buildCommonIPsQuery(userA, userB))
}

func (q *UserAltQuery) GetCommonApps(userA, userB string) UserAppCounts {
	return q.getAppCounts(q.buildCommonAppsQuery(userA, userB))
}

func (q *UserAltQuery) GetDiffIPs(userA, userB string) UserIPCounts {
	return q.getIPCounts(q.buildDiffIPsQuery(userA, userB))
}

func (q *UserAltQuery) GetDiffApps(userA, userB string) UserAppCounts {
	return q.getAppCounts(q.buildDiffAppsQuery(userA, userB))
}

func (q *UserAltQuery) GetUserIPs(user string) UserIPCounts {
	return q.getIPCounts(q.buildUserIPsQuery(user))
}

func (q *UserAltQuery) GetUserApps(user string) UserAppCounts {
	return q.getAppCounts(q.buildUserAppsQuery(user))
}

func (q *UserAltQuery) baseAltQuery(user string) *sqlf.Stmt {
	return sqlf.Select("ul.name, COUNT(*) AS cnt").
		From("user_log AS ul").
		Where("ul.name <> lower(?)", user).
		Where("ol.name = lower(?)", user).
		Where("ul.first <> ol.first").
		GroupBy("ul.name").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

func (q *UserAltQuery) buildIPQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

func (q *UserAltQuery) buildDeviceQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.device = ol.device").
		Where("ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

func (q *UserAltQuery) buildAppQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("ol.app <> 0 AND ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

func (q *UserAltQuery) buildUID2Query(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.uid2 = ol.uid2").
		Where("ol.uid2 <> 0")
}

func (q *UserAltQuery) buildEmailQuery(user, email string) *sqlf.Stmt {
	emailSubQuery := sqlf.Select("id, name").
		Select("to_search_string(?) <<-> to_search_string(email) AS dist", email).
		From("users").
		OrderBy("dist ASC, id DESC").
		Limit(q.limit + 1)

	return sqlf.Select("name").
		Select("100 - round(dist * 100)").
		From("").SubQuery("(", ") AS u", emailSubQuery).
		Where("lower(name) <> lower(?)", user).
		Where("dist < 0.5")
}

func (q *UserAltQuery) buildUserIPsQuery(user string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at), ul.ip").
		From("user_log AS ul").
		Where("ul.name = lower(?)", user).
		GroupBy("ul.ip").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

func (q *UserAltQuery) buildUserAppsQuery(user string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at), ul.user_agent").
		From("user_log AS ul").
		Where("ul.name = lower(?)", user).
		GroupBy("ul.user_agent").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

func (q *UserAltQuery) baseCommonQuery(userA, userB string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at)").
		From("user_log AS ul").
		Where("ul.name = lower(?)", userA).
		Where("ol.name = lower(?)", userB).
		Where("ul.first <> ol.first").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

func (q *UserAltQuery) buildCommonIPsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseCommonQuery(userA, userB).
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.ip").
		GroupBy("ul.ip")
}

func (q *UserAltQuery) buildCommonAppsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseCommonQuery(userA, userB).
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.user_agent").
		GroupBy("ul.user_agent")
}

func (q *UserAltQuery) baseDiffQuery(userA, userB string) *sqlf.Stmt {
	sub := sqlf.Select("*").
		From("user_log").
		Where("name = ?", userB)

	return sqlf.Select("COUNT(*) AS cnt, MIN(ol.at), MAX(ol.at)").
		From("").SubQuery("(", ") AS ul", sub).
		Where("ol.name = lower(?)", userA).
		OrderBy("cnt DESC").
		Limit(q.limit)
}

func (q *UserAltQuery) buildDiffIPsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseDiffQuery(userA, userB).
		RightJoin("user_log AS ol", "ul.ip = ol.ip").
		Where("ul.ip IS NULL").
		Select("ol.ip").
		GroupBy("ol.ip")
}

func (q *UserAltQuery) buildDiffAppsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseDiffQuery(userA, userB).
		RightJoin("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("ul.app IS NULL").
		Select("ol.user_agent").
		GroupBy("ol.user_agent")
}

func (q *UserAltQuery) getAlts(stmt *sqlf.Stmt) UserAltCounts {
	q.tx.QueryStmt(stmt)
	var alts UserAltCounts

	for {
		var count UserAltCount
		ok := q.tx.Scan(&count.Alt, &count.Count)
		if !ok {
			break
		}
		count.baseURL = q.baseURL
		alts = append(alts, count)
	}

	return alts
}

func (q *UserAltQuery) getIPCounts(stmt *sqlf.Stmt) UserIPCounts {
	q.tx.QueryStmt(stmt)
	var counts UserIPCounts

	for {
		var count UserIPCount
		ok := q.tx.Scan(&count.Count, &count.From, &count.To, &count.IP)
		if !ok {
			break
		}
		count.ipAPI = q.ipAPI
		counts = append(counts, count)
	}

	return counts
}

func (q *UserAltQuery) getAppCounts(stmt *sqlf.Stmt) UserAppCounts {
	q.tx.QueryStmt(stmt)
	var counts UserAppCounts

	for {
		var count UserAppCount
		ok := q.tx.Scan(&count.Count, &count.From, &count.To, &count.UserAgent)
		if !ok {
			break
		}
		counts = append(counts, count)
	}

	return counts
}

func (c UserAltCount) String() string {
	const userAltFormat = `<a href="%susers/%s">%s</a> (%d)`
	return fmt.Sprintf(userAltFormat, c.baseURL, c.Alt, c.Alt, c.Count)
}

func (c UserIPCount) String() string {
	formattedIP := fmt.Sprintf(c.ipAPI, c.IP, c.IP)
	return fmt.Sprintf("%s (%d, %s — %s)",
		formattedIP,
		c.Count,
		c.From.Format("02.01.2006"),
		c.To.Format("02.01.2006"))
}

func (c UserAppCount) String() string {
	uaData := useragent.Parse(c.UserAgent)
	formattedApp := c.UserAgent
	if !uaData.IsUnknown() {
		formattedApp = fmt.Sprintf("%s %s on %s", uaData.Name, uaData.Version, uaData.OS)
		if uaData.OSVersion != "" {
			formattedApp += " " + uaData.OSVersion
		}
		if uaData.Device != "" {
			formattedApp += ", " + uaData.Device
		}
	}

	return fmt.Sprintf("%s (%d, %s — %s)",
		formattedApp,
		c.Count,
		c.From.Format("02.01.2006"),
		c.To.Format("02.01.2006"))
}

func (counts UserAltCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, ", ")
}

func (counts UserIPCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, "\n")
}

func (counts UserAppCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, "\n")
}
