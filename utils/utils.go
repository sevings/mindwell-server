package utils

import (
	"bytes"
	crypto "crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
	goconf "github.com/zpatrick/go-config"

	"github.com/go-openapi/errors"
	"github.com/sevings/mindwell-server/models"

	// to use postgres
	_ "github.com/lib/pq"
)

var errUnauthorized = errors.New(401, "Invalid or expired API key")
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

// CanViewEntry returns true if the user is allowed to read the entry.
func CanViewEntry(tx *AutoTx, userID, entryID int64) bool {
	return tx.QueryBool("SELECT can_view_entry($1, $2)", userID, entryID)
}

// IsOpenForMe returns true if you can see \param name tlog
func IsOpenForMe(tx *AutoTx, userID *models.UserID, name interface{}) bool {
	return tx.QueryBool("SELECT can_view_tlog($1, $2)", userID.ID, name)
}

const userIDQuery = `
			SELECT id, name, followers_count, 
				invited_by is not null, karma < -1, verified,
				invite_ban > CURRENT_DATE, vote_ban > CURRENT_DATE, 
				comment_ban > CURRENT_DATE, live_ban > CURRENT_DATE
			FROM users `

func LoadUserIDByID(tx *AutoTx, id int64) (*models.UserID, error) {
	const q = userIDQuery + "WHERE id = $1"
	tx.Query(q, id)
	return scanUserID(tx)
}

func LoadUserIDByApiKey(tx *AutoTx, apiKey string) (*models.UserID, error) {
	const q = userIDQuery + "WHERE api_key = $1 AND valid_thru > CURRENT_TIMESTAMP"
	tx.Query(q, apiKey)
	return scanUserID(tx)
}

func scanUserID(tx *AutoTx) (*models.UserID, error) {
	var user models.UserID
	user.Ban = &models.UserIDBan{}
	tx.Scan(&user.ID, &user.Name, &user.FollowersCount,
		&user.IsInvited, &user.NegKarma, &user.Verified,
		&user.Ban.Invite, &user.Ban.Vote, &user.Ban.Comment, &user.Ban.Live)
	if tx.Error() != nil {
		return nil, errUnauthorized
	}

	user.Ban.Invite = user.Ban.Invite || !user.IsInvited || !user.Verified
	user.Ban.Vote = user.Ban.Vote || !user.IsInvited || user.NegKarma || !user.Verified
	user.Ban.Comment = user.Ban.Comment || !user.IsInvited || !user.Verified
	user.Ban.Live = user.Ban.Live || !user.Verified

	return &user, nil
}

func readUserID(secret []byte, tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		log.Println(err)
		return 0, errUnauthorized
	}

	if !token.Valid {
		log.Printf("Invalid token: %s\n", tokenString)
		return 0, errUnauthorized

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Error get claims: %s\n", tokenString)
		return 0, errUnauthorized
	}

	if claims.Valid() != nil {
		return 0, errUnauthorized
	}

	id, err := strconv.ParseInt(claims["sub"].(string), 32, 64)
	if err != nil {
		log.Println(err)
		return 0, errUnauthorized
	}

	return id, nil
}

func NewKeyAuth(db *sql.DB, secret []byte) func(apiKey string) (*models.UserID, error) {
	return func(apiKey string) (*models.UserID, error) {
		tx := NewAutoTx(db)
		defer tx.Finish()

		if len(apiKey) == 32 {
			return LoadUserIDByApiKey(tx, apiKey)
		}

		id, err := readUserID(secret, apiKey)
		if err != nil {
			return nil, err
		}

		return LoadUserIDByID(tx, id)
	}
}

func BuildApiToken(secret []byte, userID int64) (string, int64) {
	now := time.Now().Unix()
	exp := now + 60*60*24*365
	sub := strconv.FormatInt(userID, 32)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": now,
		"exp": exp,
		"sub": sub,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Print(err)
	}

	return tokenString, exp
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
