package userutil

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/leporo/sqlf"
	"github.com/mileusna/useragent"
	"github.com/sevings/mindwell-server/lib/database"
)

// UserAltQuery provides methods for detecting alternative user accounts.
type UserAltQuery struct {
	tx      *database.AutoTx
	limit   int
	baseURL string
	ipAPI   string
}

// UserAltCount represents a count of alternative account matches.
type UserAltCount struct {
	Count   int64
	Alt     string
	baseURL string
}

// UserIPCount represents IP address usage statistics.
type UserIPCount struct {
	Count int64
	From  time.Time
	To    time.Time
	IP    string
	ipAPI string
}

// UserAppCount represents user agent usage statistics.
type UserAppCount struct {
	Count     int64
	From      time.Time
	To        time.Time
	UserAgent string
}

// UserAltCounts is a collection of alternative account counts.
type UserAltCounts []UserAltCount

// UserIPCounts is a collection of IP address counts.
type UserIPCounts []UserIPCount

// UserAppCounts is a collection of user agent counts.
type UserAppCounts []UserAppCount

// NewUserAltQuery creates a new UserAltQuery instance.
func NewUserAltQuery(tx *database.AutoTx, limit int, baseURL, ipAPI string) *UserAltQuery {
	return &UserAltQuery{
		tx:      tx,
		limit:   limit,
		baseURL: baseURL,
		ipAPI:   ipAPI,
	}
}

// GetSuspectedAlt returns the suspected alternative account for a user.
func (q *UserAltQuery) GetSuspectedAlt(user string) (string, bool) {
	var alt sql.NullString
	var conf bool
	confQuery := sqlf.Select("alt_of, confirmed_alt").
		From("users").
		Where("lower(name) = lower(?)", user)
	q.tx.QueryStmt(confQuery).Scan(&alt, &conf)
	return alt.String, conf
}

// GetIPAlts finds alternative accounts based on shared IP addresses.
func (q *UserAltQuery) GetIPAlts(user string) UserAltCounts {
	return q.getAlts(q.buildIPQuery(user))
}

// GetDeviceAlts finds alternative accounts based on shared devices.
func (q *UserAltQuery) GetDeviceAlts(user string) UserAltCounts {
	return q.getAlts(q.buildDeviceQuery(user))
}

// GetAppAlts finds alternative accounts based on shared app/device combinations.
func (q *UserAltQuery) GetAppAlts(user string) UserAltCounts {
	return q.getAlts(q.buildAppQuery(user))
}

// GetUID2Alts finds alternative accounts based on shared UID2 tracking.
func (q *UserAltQuery) GetUID2Alts(user string) UserAltCounts {
	return q.getAlts(q.buildUID2Query(user))
}

// GetEmailAlts finds alternative accounts based on email similarity.
func (q *UserAltQuery) GetEmailAlts(user, email string) UserAltCounts {
	return q.getAlts(q.buildEmailQuery(user, email))
}

// GetCommonIPs finds IP addresses shared between two users.
func (q *UserAltQuery) GetCommonIPs(userA, userB string) UserIPCounts {
	return q.getIPCounts(q.buildCommonIPsQuery(userA, userB))
}

// GetCommonApps finds user agents shared between two users.
func (q *UserAltQuery) GetCommonApps(userA, userB string) UserAppCounts {
	return q.getAppCounts(q.buildCommonAppsQuery(userA, userB))
}

// GetDiffIPs finds IP addresses used by userA but not userB.
func (q *UserAltQuery) GetDiffIPs(userA, userB string) UserIPCounts {
	return q.getIPCounts(q.buildDiffIPsQuery(userA, userB))
}

// GetDiffApps finds user agents used by userA but not userB.
func (q *UserAltQuery) GetDiffApps(userA, userB string) UserAppCounts {
	return q.getAppCounts(q.buildDiffAppsQuery(userA, userB))
}

// GetUserIPs gets all IP addresses used by a user.
func (q *UserAltQuery) GetUserIPs(user string) UserIPCounts {
	return q.getIPCounts(q.buildUserIPsQuery(user))
}

// GetUserApps gets all user agents used by a user.
func (q *UserAltQuery) GetUserApps(user string) UserAppCounts {
	return q.getAppCounts(q.buildUserAppsQuery(user))
}

// baseAltQuery creates the base query for finding alternative accounts.
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

// buildIPQuery builds query for IP-based alternative account detection.
func (q *UserAltQuery) buildIPQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

// buildDeviceQuery builds query for device-based alternative account detection.
func (q *UserAltQuery) buildDeviceQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.device = ol.device").
		Where("ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

// buildAppQuery builds query for app/device-based alternative account detection.
func (q *UserAltQuery) buildAppQuery(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("ol.app <> 0 AND ol.device <> 0").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'")
}

// buildUID2Query builds query for UID2-based alternative account detection.
func (q *UserAltQuery) buildUID2Query(user string) *sqlf.Stmt {
	return q.baseAltQuery(user).
		Join("user_log AS ol", "ul.uid2 = ol.uid2").
		Where("ol.uid2 <> 0")
}

// buildEmailQuery builds query for email similarity-based alternative account detection.
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

// buildUserIPsQuery builds query for getting all IPs used by a user.
func (q *UserAltQuery) buildUserIPsQuery(user string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at), ul.ip").
		From("user_log AS ul").
		Where("ul.name = lower(?)", user).
		GroupBy("ul.ip").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

// buildUserAppsQuery builds query for getting all user agents used by a user.
func (q *UserAltQuery) buildUserAppsQuery(user string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at), ul.user_agent").
		From("user_log AS ul").
		Where("ul.name = lower(?)", user).
		GroupBy("ul.user_agent").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

// baseCommonQuery creates the base query for finding common attributes between users.
func (q *UserAltQuery) baseCommonQuery(userA, userB string) *sqlf.Stmt {
	return sqlf.Select("COUNT(*) AS cnt, MIN(ul.at), MAX(ul.at)").
		From("user_log AS ul").
		Where("ul.name = lower(?)", userA).
		Where("ol.name = lower(?)", userB).
		Where("ul.first <> ol.first").
		OrderBy("cnt DESC").
		Limit(q.limit)
}

// buildCommonIPsQuery builds query for finding shared IPs between two users.
func (q *UserAltQuery) buildCommonIPsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseCommonQuery(userA, userB).
		Join("user_log AS ol", "ul.ip = ol.ip").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.ip").
		GroupBy("ul.ip")
}

// buildCommonAppsQuery builds query for finding shared user agents between two users.
func (q *UserAltQuery) buildCommonAppsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseCommonQuery(userA, userB).
		Join("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("(CASE WHEN ul.at > ol.at THEN ul.at - ol.at ELSE ol.at - ul.at END) < interval '1 hour'").
		Select("ul.user_agent").
		GroupBy("ul.user_agent")
}

// baseDiffQuery creates the base query for finding different attributes between users.
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

// buildDiffIPsQuery builds query for finding IPs unique to userA.
func (q *UserAltQuery) buildDiffIPsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseDiffQuery(userA, userB).
		RightJoin("user_log AS ol", "ul.ip = ol.ip").
		Where("ul.ip IS NULL").
		Select("ol.ip").
		GroupBy("ol.ip")
}

// buildDiffAppsQuery builds query for finding user agents unique to userA.
func (q *UserAltQuery) buildDiffAppsQuery(userA, userB string) *sqlf.Stmt {
	return q.baseDiffQuery(userA, userB).
		RightJoin("user_log AS ol", "ul.app = ol.app AND ul.device = ol.device").
		Where("ul.app IS NULL").
		Select("ol.user_agent").
		GroupBy("ol.user_agent")
}

// getAlts executes an alt detection query and returns the results.
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

// getIPCounts executes an IP count query and returns the results.
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

// getAppCounts executes a user agent count query and returns the results.
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

// String formats UserAltCount as an HTML link.
func (c UserAltCount) String() string {
	const userAltFormat = `<a href="%susers/%s">%s</a> (%d)`
	return fmt.Sprintf(userAltFormat, c.baseURL, c.Alt, c.Alt, c.Count)
}

// String formats UserIPCount with IP API link and date range.
func (c UserIPCount) String() string {
	formattedIP := fmt.Sprintf(c.ipAPI, c.IP, c.IP)
	return fmt.Sprintf("%s (%d, %s — %s)",
		formattedIP,
		c.Count,
		c.From.Format("02.01.2006"),
		c.To.Format("02.01.2006"))
}

// String formats UserAppCount with parsed user agent information.
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

// String formats UserAltCounts as a comma-separated list.
func (counts UserAltCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, ", ")
}

// String formats UserIPCounts as a newline-separated list.
func (counts UserIPCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, "\n")
}

// String formats UserAppCounts as a newline-separated list.
func (counts UserAppCounts) String() string {
	var strs []string
	for _, c := range counts {
		strs = append(strs, c.String())
	}
	return strings.Join(strs, "\n")
}
