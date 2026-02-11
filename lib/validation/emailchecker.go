package validation

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
)

const checkUrl = "http://open.kickbox.com/v1/disposable/"

type reply struct {
	Disposable bool
}

// EmailAllowedChecker is the interface for email validation
type EmailAllowedChecker interface {
	IsAllowed(email string) bool
}

// ConfigProvider provides access to configuration strings
type ConfigProvider interface {
	ConfigStrings(key string) []string
	LogSystem() *zap.Logger
}

type EmailChecker struct {
	srv     ConfigProvider
	client  *http.Client
	trusted []string
	banned  []string
}

func NewEmailChecker(srv ConfigProvider) *EmailChecker {
	return &EmailChecker{
		srv: srv,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		trusted: srv.ConfigStrings("server.trust_email"),
		banned:  srv.ConfigStrings("server.ban_email"),
	}
}

func (ec *EmailChecker) IsAllowed(email string) bool {
	loginAtService := strings.Split(email, "@")
	if len(loginAtService) < 2 {
		return false
	}

	service := loginAtService[1]

	if slices.Contains(ec.trusted, service) {
		return true
	}

	if slices.Contains(ec.banned, service) {
		return false
	}

	resp, err := ec.client.Get(checkUrl + service)
	if err != nil {
		ec.logError(err)
		return true
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ec.logError(err)
		return true
	}

	if resp.StatusCode != 200 {
		ec.srv.LogSystem().Error(string(body))
		return true
	}

	var result reply
	err = json.Unmarshal(body, &result)
	if err != nil {
		ec.logError(err)
		return true
	}

	if result.Disposable {
		ec.banned = append(ec.banned, service)
	}

	return !result.Disposable
}

func (ec *EmailChecker) logError(err error) {
	ec.srv.LogSystem().Error(err.Error())
}
