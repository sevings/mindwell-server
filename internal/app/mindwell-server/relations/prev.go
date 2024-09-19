package relations

import (
	"github.com/patrickmn/go-cache"
	"github.com/sevings/mindwell-server/models"
	"time"
)

var prevFollowings *cache.Cache

func init() {
	prevFollowings = cache.New(time.Hour, time.Hour)
}

func checkPrev(userID *models.UserID, toName string) (found bool) {
	_, found = prevFollowings.Get(userID.Name + ":" + toName)
	return
}

func setPrev(userID *models.UserID, toName string, relation *models.Relationship) {
	prevFollowings.SetDefault(userID.Name+":"+toName, relation)
}

func removePrev(userID *models.UserID, toName string) {
	prevFollowings.Delete(userID.Name + ":" + toName)
}
