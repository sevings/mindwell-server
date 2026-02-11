package textutil

import (
	"bytes"
	crypto "crypto/rand"
	"encoding/base64"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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

func CutHtml(text string, maxLineCount, maxLineLen, maxImageCount int) (string, bool) {
	var tags []string
	var hasTags bool
	idx := 0
	lineCount := 0
	lineLen := 0.
	imageCount := 0
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

		if tag == "img" {
			if imageCount >= maxImageCount || lineCount >= maxLineCount-1 {
				idx = tagStart - 1
				break
			}

			imageCount++
		}

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
		} else if tag != "br" && tag != "img" {
			tags = append(tags, tag)
			hasTags = true

			k, ok := fontKoeff(tag)
			if ok {
				koeff = k
			}
		}

		for idx < textLen {
			tagChar := text[idx]
			idx++
			if tagChar == '>' {
				break
			}
		}

		if tag == "br" || (isClosing && tag == "p") || (isClosing && tag == "blockquote") {
			lineCount++
			lineLen = 0
		} else if tag == "img" {
			lineCount += 2
			lineLen = 0
		}
	}

	if idx >= textLen {
		return text, false
	}

	resLen := idx
	c, _ := utf8.DecodeRuneInString(text[idx:])
	wasSpace := unicode.IsSpace(c)
	var wasImg bool
	for idx > 0 {
		c, size := utf8.DecodeLastRune([]byte(text)[:idx])

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
			} else if len(tag) > 4 && tag[1:4] == "img" {
				wasImg = true
				break
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
		resLen--
	} else {
		resLen = idx
	}
	text = text[:resLen]

	if !wasImg && hasTags && len(tags) == 0 {
		return text + "<p>…</p>", true
	}
	if !wasImg && !hasTags {
		return text + "…", true
	}
	if len(tags) == 0 {
		return text, true
	}

	var res bytes.Buffer
	res.Grow(resLen + len(tags)*4)
	res.WriteString(text)
	if !wasImg {
		res.WriteString("…")
	}

	for i := len(tags) - 1; i >= 0; i-- {
		res.WriteString("</")
		res.WriteString(tags[i])
		res.WriteString(">")
	}

	return res.String(), true
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
