package helper

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"github.com/sevings/mindwell-server/utils"
	"hash/adler32"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type userRequest struct {
	User string  `json:"user"`
	Ip   string  `json:"ip"`
	Ua   string  `json:"user_agent"`
	Dev  string  `json:"dev"`
	App  string  `json:"browser"`
	At   float64 `json:"ts"`
	Type string  `json:"type"`
	Mtd  string  `json:"method"`
	Url  string  `json:"url"`
}

func (req userRequest) key() string {
	var str strings.Builder
	str.Grow(100)

	str.WriteString(req.Ip)
	str.WriteString(req.Dev)
	str.WriteString(req.App)
	str.WriteString(req.User)

	return str.String()
}

func ImportUserLog(tx *utils.AutoTx) {
	prev := make(map[string]userRequest)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "{") {
			continue
		}

		var req userRequest
		err := json.Unmarshal(scanner.Bytes(), &req)
		if err != nil {
			log.Println(err)
			continue
		}

		if req.Type != "web" || req.Mtd != "GET" ||
			(req.Url != "/live" && req.Url != "/friends" && req.Url != "/best" && req.Url != "/users") {
			continue
		}

		if req.User == "" || req.Dev == "" || req.App == "" {
			continue
		}

		prev[req.key()] = req
	}

	for _, req := range prev {
		saveUserRequest(tx, req)
	}

	log.Printf("Added %d log lines.\n", len(prev))
}

func saveUserRequest(tx *utils.AutoTx, req userRequest) {
	app, err := strconv.ParseUint(req.App[:16], 16, 64)
	if err != nil {
		log.Printf("Browser id is invalid: %s\n", req.App)
	}

	device, err := base64.StdEncoding.DecodeString(req.Dev)
	if err != nil {
		log.Printf("Device id is invalid: %s\n", req.Dev)
	}
	for _, c := range []byte(";-180") {
		device = append(device, c)
	}
	dev := adler32.Checksum(device)

	at := time.UnixMicro(int64(req.At * 1000000))
	ip := strings.SplitN(req.Ip, ",", 2)[0]

	const query = `
    INSERT INTO user_log(name, ip, user_agent, device, app, uid, at, first) 
    VALUES(lower($1), $2, $3, $4, $5, $6, $7, $8)
`

	tx.Exec(query, req.User, ip, req.Ua, int32(dev), int64(app), -1, at, false)
}
