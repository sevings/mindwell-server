package utils

import (
	"bytes"
	crypto "crypto/rand"
	"encoding/base64"
	"github.com/leporo/sqlf"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	goconf "github.com/zpatrick/go-config"

	"github.com/sevings/mindwell-server/models"

	// to use postgres
	_ "github.com/lib/pq"
)

var htmlEsc = strings.NewReplacer(
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;",
	"\"", "&#34;",
	"'", "&#39;",
	"\n", "<br>",
	"\r", "",
)

// LoadConfig creates app config from file
func LoadConfig(fileName string) *goconf.Config {
	toml := goconf.NewTOMLFile(fileName + ".toml")
	loader := goconf.NewOnceLoader(toml)
	config := goconf.NewConfig([]goconf.Provider{loader})
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}

	return config
}

func AddAuthorRelationsToMeQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	return q.
		With("relations_to_me",
			sqlf.Select("relation.type, relations.from_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.to_id = ?", userID.ID)).
		LeftJoin("relations_to_me", "relations_to_me.from_id = authors.id")
}

func AddAuthorRelationsFromMeQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	return q.
		With("relations_from_me",
			sqlf.Select("relation.type, relations.to_id").
				From("relations").
				Join("relation", "relations.type = relation.id").
				Where("relations.from_id = ?", userID.ID)).
		LeftJoin("relations_from_me", "relations_from_me.to_id = authors.id")
}

func AddEntryOpenQuery(q *sqlf.Stmt, userID *models.UserID, showMe bool) *sqlf.Stmt {
	return q.Where(`
CASE entry_privacy.type
WHEN 'all' THEN TRUE
WHEN 'registered' THEN ?
WHEN 'invited' THEN ?
WHEN 'followers' THEN authors.id = ? OR authors.creator_id = ? OR relations_from_me.type = 'followed'
WHEN 'some' THEN authors.id = ? OR authors.creator_id = ?
	OR EXISTS(SELECT 1 from entries_privacy WHERE user_id = ? AND entry_id = entries.id)
WHEN 'me' THEN ? AND authors.id = ?
ELSE FALSE
END
`, userID.ID > 0, userID.IsInvited, userID.ID, userID.ID, userID.ID, userID.ID, userID.ID, showMe, userID.ID)
}

func AddCanViewAuthorQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddAuthorRelationsFromMeQuery(q, userID)
	AddAuthorRelationsToMeQuery(q, userID)

	return q.
		Where("(relations_to_me.type IS NULL OR relations_to_me.type <> 'ignored')").
		Where("(relations_from_me.type IS NULL OR relations_from_me.type <> 'ignored')").
		Where(`
CASE user_privacy.type
WHEN 'all' THEN TRUE
WHEN 'registered' THEN ?
WHEN 'invited' THEN ?
WHEN 'followers' THEN authors.id = ? OR relations_from_me.type = 'followed'
ELSE FALSE
END`, userID.ID > 0, userID.IsInvited, userID.ID)
}

func AddCanViewEntryQuery(q *sqlf.Stmt, userID *models.UserID) *sqlf.Stmt {
	AddEntryOpenQuery(q, userID, true)
	return AddCanViewAuthorQuery(q, userID)
}

func LoadRelation(tx *AutoTx, from, to int64) string {
	if from == 0 || to == 0 {
		return models.RelationshipRelationNone
	}

	relationQuery := sqlf.Select("relation.type").
		From("relations").
		Join("relation", "relation.id = relations.type").
		Where("from_id = ?", from).
		Where("to_id = ?", to)

	var relation string
	tx.QueryStmt(relationQuery).Scan(&relation)

	if relation == "" {
		return models.RelationshipRelationNone
	}

	return relation
}

func CanViewTlogName(tx *AutoTx, userID *models.UserID, tlog string) bool {
	if tlog == "" {
		return false
	}

	q := sqlf.Select("TRUE").
		From("users AS authors").
		LeftJoin("user_privacy", "authors.privacy = user_privacy.id").
		Where("lower(name) = lower(?)", tlog)

	AddCanViewAuthorQuery(q, userID)
	return tx.QueryStmt(q).ScanBool()
}

// CanViewEntry returns true if the user is allowed to read the entry.
func CanViewEntry(tx *AutoTx, userID *models.UserID, entryID int64) bool {
	if entryID == 0 {
		return false
	}

	q := sqlf.Select("TRUE").
		From("entries").
		Join("entry_privacy", "entries.visible_for = entry_privacy.id").
		Join("users AS authors", "entries.author_id = authors.id").
		Join("user_privacy", "authors.privacy = user_privacy.id").
		Where("entries.id = ?", entryID)

	AddCanViewEntryQuery(q, userID)

	return tx.QueryStmt(q).ScanBool()
}

func queryUser(tx *AutoTx, query string, arg any) *models.User {
	var user models.User

	tx.Query(query, arg).Scan(&user.ID, &user.Name, &user.ShowName,
		&user.IsOnline, &user.IsTheme)

	return &user
}

func LoadUser(tx *AutoTx, id int64) *models.User {
	const query = `
SELECT id, name, show_name,
is_online(last_seen_at) AND creator_id IS NULL, creator_id IS NOT NULL
FROM users
WHERE id = $1`

	return queryUser(tx, query, id)
}

func LoadUserByName(tx *AutoTx, name string) *models.User {
	const query = `
SELECT id, name, show_name,
is_online(last_seen_at) AND creator_id IS NULL, creator_id IS NOT NULL
FROM users
WHERE lower(name) = lower($1)`

	return queryUser(tx, query, name)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GenerateString returns random string
func GenerateString(length int) string {
	bytesLen := length*6/8 + 1
	data := make([]byte, bytesLen)
	_, err := crypto.Read(data)
	if err == nil {
		str := base64.URLEncoding.EncodeToString(data)
		return str[:length]
	}

	log.Print(err)

	// fallback on error
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := len(b)-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var tagRe = regexp.MustCompile(`<[^>]+>`)
var imgRe = regexp.MustCompile(`(?i)<img[^>]+>`)

func RemoveHTML(text string) string {
	text = imgRe.ReplaceAllString(text, "[изображение]")
	text = tagRe.ReplaceAllString(text, "")

	return text
}

func CutText(text string, size int) (string, bool) {
	runes := []rune(text)

	if len(runes) <= size {
		return text, false
	}

	text = string(runes[:size])

	var isSpace bool
	trim := func(r rune) bool {
		if isSpace {
			return unicode.IsSpace(r)
		}

		isSpace = unicode.IsSpace(r)
		return true
	}
	text = strings.TrimRightFunc(text, trim)

	text += "..."

	return text, true
}

func CutHtml(text string, maxLineCount, maxLineLen int) (string, bool) {
	var output bytes.Buffer
	var tags []string
	idx := 0
	lineCount := 0
	lineLen := 0.
	koeff := 1.
	textLen := len(text)

	fontKoeff := func(tag string) (float64, bool) {
		switch tag {
		case "h1":
			return 8.16, true
		case "h2":
			return 5.22, true
		case "h3":
			return 4, true
		case "h4":
			return 2.47, true
		case "h5":
			return 1.31, true
		case "h6":
			return 1, true
		case "blockquote":
			return 4.61, true
		default:
			return 0, false
		}
	}

	for idx < textLen && lineCount < maxLineCount {
		c, s := utf8.DecodeRuneInString(text[idx:])
		if c == utf8.RuneError {
			idx += s
			continue
		}

		idx += s
		output.WriteRune(c)

		if c != '<' {
			lineLen += koeff
			if lineLen > float64(maxLineLen) {
				lineLen = koeff
				lineCount++
			}

			continue
		}

		if idx >= textLen-1 {
			break
		}

		isClosing := text[idx] == '/'
		if isClosing {
			idx++
			output.WriteByte('/')
		}

		tagStart := idx
		for idx < textLen {
			tagChar := text[idx]
			if tagChar != ' ' && tagChar != '>' {
				idx++
			} else {
				break
			}
		}

		tag := text[tagStart:idx]
		output.WriteString(tag)

		if isClosing {
			for i := len(tags) - 1; i >= 0; i-- {
				if tags[i] != tag {
					continue
				}

				tags = tags[:i]
				break
			}

			_, ok := fontKoeff(tag)
			if ok {
				koeff = 1
			}
		} else if tag != "br" {
			tags = append(tags, tag)

			k, ok := fontKoeff(tag)
			if ok {
				koeff = k
			}
		}

		for idx < textLen {
			tagChar := text[idx]
			if tagChar == '>' {
				break
			}

			idx++
			output.WriteByte(tagChar)
		}

		idx++
		output.WriteByte('>')

		if tag == "br" || (isClosing && tag == "p") || (isClosing && tag == "blockquote") {
			lineCount++
			lineLen = 0
		}
	}

	if idx >= textLen {
		return text, false
	}

	c, _ := utf8.DecodeRuneInString(text[idx:])
	wasSpace := unicode.IsSpace(c)
	for idx = output.Len(); idx > 0; {
		c, size := utf8.DecodeLastRune(output.Bytes()[:idx])

		isSpace := unicode.IsSpace(c) || unicode.IsPunct(c)
		if c == '>' {
			tagStart := idx - 1
			for ; tagStart >= 0; tagStart-- {
				if text[tagStart] == '<' {
					break
				}
			}

			tag := text[tagStart:idx]
			if tag == "</p>" {
				tags = append(tags, "p")
				isSpace = true
				size = 4
			} else if tag[1] != '/' {
				isSpace = true
				size = len(tag)
			}
		}

		if wasSpace {
			if isSpace {
				idx -= size
				continue
			}

			break
		}

		idx -= size
		wasSpace = isSpace
	}

	if idx == 0 {
		idx = output.Len() - 1
	}

	output.Truncate(idx)
	output.WriteString("…")

	for i := len(tags) - 1; i >= 0; i-- {
		output.WriteString("</")
		output.WriteString(tags[i])
		output.WriteString(">")
	}

	return output.String(), true
}

func ParseFloat(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if len(val) > 0 && err != nil {
		log.Printf("error parse float: '%s'", val)
	}

	return res
}

func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 6, 64)
}

func ParseInt64(val string) int64 {
	res, err := strconv.ParseInt(val, 32, 64)
	if len(val) > 0 && err != nil {
		log.Printf("error parse int: '%s'", val)
	}

	return res
}

func FormatInt64(val int64) string {
	return strconv.FormatInt(val, 32)
}

func ReplaceToHtml(val string) string {
	return htmlEsc.Replace(val)
}

func HideEmail(email string) string {
	nameLen := strings.LastIndex(email, "@")

	if nameLen < 0 {
		return ""
	}

	if nameLen < 3 {
		return "***" + email[nameLen:]
	}

	return email[:1] + "***" + email[nameLen-1:]
}
