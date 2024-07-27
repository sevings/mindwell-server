package test

import (
	"fmt"
	"github.com/sevings/mindwell-server/utils"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/entries"
	"github.com/sevings/mindwell-server/restapi/operations/favorites"
	"github.com/sevings/mindwell-server/restapi/operations/me"
	"github.com/sevings/mindwell-server/restapi/operations/users"
)

func BenchmarkLoadLive(b *testing.B) {
	post := api.MePostMeTlogHandler.Handle
	var title string
	votable := true
	commentable := true
	live := true
	shared := true
	draft := false
	entryParams := me.PostMeTlogParams{
		Title:         &title,
		Privacy:       models.EntryPrivacyAll,
		IsVotable:     &votable,
		IsCommentable: &commentable,
		InLive:        &live,
		IsShared:      &shared,
		IsDraft:       &draft,
	}
	for i := 0; i < 1000; i++ {
		title = fmt.Sprintf("Entry %d", i)
		entryParams.Content = fmt.Sprintf("test test test %d", i)
		post(entryParams, userIDs[0])
	}

	var limit int64 = 30
	before := "0"
	after := "0"
	section := "entries"
	query := ""
	tag := ""
	source := ""
	params := entries.GetEntriesLiveParams{
		Limit:   &limit,
		Before:  &before,
		After:   &after,
		Section: &section,
		Query:   &query,
		Tag:     &tag,
		Source:  &source,
	}

	load := api.EntriesGetEntriesLiveHandler.Handle

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		load(params, userIDs[1])
	}
}

func BenchmarkLoadFavorite(b *testing.B) {
	post := api.MePostMeTlogHandler.Handle
	var title string
	votable := true
	commentable := true
	live := true
	shared := true
	draft := false
	entryParams := me.PostMeTlogParams{
		Title:         &title,
		Privacy:       models.EntryPrivacyAll,
		IsVotable:     &votable,
		IsCommentable: &commentable,
		InLive:        &live,
		IsShared:      &shared,
		IsDraft:       &draft,
	}
	var ids []int64
	for i := 0; i < 1000; i++ {
		title = fmt.Sprintf("Entry %d", i)
		entryParams.Content = fmt.Sprintf("test test test %d", i)
		resp := post(entryParams, userIDs[0])
		body := resp.(*me.PostMeTlogCreated)
		id := body.Payload.ID
		ids = append(ids, id)
	}

	fav := api.FavoritesPutEntriesIDFavoriteHandler.Handle
	for _, id := range ids {
		favParams := favorites.PutEntriesIDFavoriteParams{ID: id}
		fav(favParams, userIDs[1])
	}

	var limit int64 = 30
	before := "0"
	after := "0"
	query := ""
	params := users.GetUsersNameFavoritesParams{
		Limit:  &limit,
		Before: &before,
		After:  &after,
		Name:   userIDs[1].Name,
		Query:  &query,
	}

	load := api.UsersGetUsersNameFavoritesHandler.Handle

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		load(params, userIDs[1])
	}
}

func TestConcurrentRequests(t *testing.T) {
	t.SkipNow()

	create := func(userID *models.UserID) {
		const baseText = `
Sentiments two *occasional affronting solicitude* travelling and one contrasted. Fortune day out married parties. Happiness remainder joy but earnestly for off. Took sold add play may none him few. If as increasing contrasted entreaties be. Now summer who day looked our behind moment coming. Pain son rose more park way that. An stairs as be lovers uneasy.

It allowance prevailed enjoyment in it. ~~Calling observe for who pressed raising his.~~ Can connection instrument astonished unaffected his motionless preference. Announcing say boy precaution unaffected difficulty alteration him. Above be would at so going heard. Engaged at village at am equally proceed. Settle nay length almost ham direct extent. Agreement for listening remainder get attention law acuteness day. Now whatever surprise resolved elegance indulged own way outlived.

***

Up am intention on dependent questions oh elsewhere september. No betrayed pleasure possible jointures we in throwing. And can event rapid any shall woman green. Hope they dear who its bred. Smiling nothing affixed he carried it clothes calling he no. Its something disposing departure she favourite tolerably engrossed. Truth short folly court why she their balls. Excellence put unaffected reasonable mrs introduced conviction she. Nay particular delightful but unpleasant for uncommonly who.

Boy favourable day can <introduced> sentiments &entreaties. Noisier carried of in warrant because. So mr plate seems cause chief widen first. Two differed husbands met screened his. Bed was form wife out ask draw. Wholly coming at we no enable. Offending sir delivered questions now new met. Acceptance she interested new boisterous day discretion celebrated.

Dependent certainty off **discovery** him his tolerably offending. Ham for attention remainder sometimes additions recommend fat our. Direction has strangers now believing. Respect enjoyed gay far exposed parlors towards. Enjoyment use tolerably dependent listening men. No peculiar in handsome together unlocked do by. Article concern joy anxious did picture sir her. Although desirous not recurred disposed off shy you numerous securing.

`
		num := rand.Int63()
		title := fmt.Sprintf("Entry %d", num)
		content := baseText + fmt.Sprintf("%d", num)
		votable := num%11 != 0
		commentable := num%13 != 0
		live := num%3 != 0
		shared := num%17 != 0
		draft := false
		tags := make([]string, 0)
		if num%7 == 0 {
			tags = append(tags, fmt.Sprintf("tag%d", num%7))
			tags = append(tags, fmt.Sprintf("tag%d", num%3))
			tags = append(tags, fmt.Sprintf("tag%d", num%2))
		}
		entryParams := me.PostMeTlogParams{
			Title:         &title,
			Content:       content,
			Tags:          tags,
			Privacy:       models.EntryPrivacyAll,
			IsVotable:     &votable,
			IsCommentable: &commentable,
			InLive:        &live,
			IsShared:      &shared,
			IsDraft:       &draft,
		}
		post := api.MePostMeTlogHandler.Handle
		resp := post(entryParams, userID)
		_, ok := resp.(*me.PostMeTlogCreated)
		require.True(t, ok)
	}

	loadLive := func(userID *models.UserID) {
		var limit int64 = 30
		before := "0"
		after := "0"
		section := "entries"
		query := ""
		tag := ""
		source := ""
		params := entries.GetEntriesLiveParams{
			Limit:   &limit,
			Before:  &before,
			After:   &after,
			Section: &section,
			Query:   &query,
			Tag:     &tag,
			Source:  &source,
		}

		load := api.EntriesGetEntriesLiveHandler.Handle
		resp := load(params, userID)
		_, ok := resp.(*entries.GetEntriesLiveOK)
		require.True(t, ok)
	}

	loadLast := func(userID *models.UserID) {
		var limit int64 = 30
		before := "0"
		after := "0"
		query := ""
		tag := ""
		sort := ""
		params := users.GetUsersNameTlogParams{
			Name:   userIDs[1].Name,
			Limit:  &limit,
			Before: &before,
			After:  &after,
			Tag:    &tag,
			Sort:   &sort,
			Query:  &query,
		}

		loadTlog := api.UsersGetUsersNameTlogHandler.Handle
		tlog := loadTlog(params, userID)
		body, ok := tlog.(*users.GetUsersNameTlogOK)
		require.True(t, ok)
		if !ok {
			return
		}

		feed := body.Payload
		if len(feed.Entries) == 0 {
			return
		}

		lastID := feed.Entries[0].ID

		load := api.EntriesGetEntriesIDHandler.Handle
		resp := load(entries.GetEntriesIDParams{ID: lastID}, userID)
		_, ok = resp.(*entries.GetEntriesIDOK)
		require.True(t, ok)
	}

	var wg sync.WaitGroup
	timeout := 10 * time.Second
	counts := make([]int, 10)

	runFunc := func(n int, f func(id *models.UserID), userID *models.UserID) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			timer := time.NewTimer(timeout)
			cnt := 0
			for {
				select {
				case <-timer.C:
					counts[n] = cnt
					return
				default:
					f(userID)
					cnt++
				}
			}
		}()
	}

	runFunc(0, create, userIDs[0])
	runFunc(1, create, userIDs[1])
	runFunc(2, create, userIDs[2])
	runFunc(3, loadLive, userIDs[0])
	runFunc(4, loadLive, userIDs[1])
	runFunc(5, loadLive, userIDs[2])
	runFunc(6, loadLive, userIDs[3])
	runFunc(7, loadLive, utils.NoAuthUser())
	runFunc(8, loadLast, userIDs[0])
	runFunc(9, loadLast, utils.NoAuthUser())

	wg.Wait()

	for i, count := range counts {
		t.Logf("Func %d: avg %.2f RPS",
			i, float64(count)/timeout.Seconds())
	}
}
