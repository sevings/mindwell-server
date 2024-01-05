package test

import (
	"fmt"
	wishesImpl "github.com/sevings/mindwell-server/internal/app/mindwell-server/wishes"
	"github.com/sevings/mindwell-server/models"
	"github.com/sevings/mindwell-server/restapi/operations/wishes"
	"github.com/stretchr/testify/require"
	"testing"
)

func checkLoadWish(t *testing.T, user *models.UserID, wishID int64, success bool) *models.Wish {
	params := wishes.GetWishesIDParams{
		ID: wishID,
	}
	load := api.WishesGetWishesIDHandler.Handle
	resp := load(params, user)
	body, ok := resp.(*wishes.GetWishesIDOK)
	require.Equal(t, success, ok)
	if !ok {
		return nil
	}

	wish := body.Payload
	require.Equal(t, wishID, wish.ID)

	if wish.State == models.WishStateNew || wish.State == models.WishStateDeclined {
		require.Empty(t, wish.Content)
		require.NotNil(t, wish.Receiver)
	} else {
		require.Equal(t, fmt.Sprintf("wish you a %d", wishID), wish.Content)
	}

	if wish.State == models.WishStateNew {
		require.NotZero(t, wish.SendUntil)
	} else {
		require.Zero(t, wish.SendUntil)
	}

	if wish.Receiver != nil {
		require.NotEqual(t, user.ID, wish.Receiver.ID)
	}

	return wish
}

func checkDeclineWish(t *testing.T, user *models.UserID, wishID int64, success bool) {
	params := wishes.DeleteWishesIDParams{ID: wishID}
	decl := api.WishesDeleteWishesIDHandler.Handle
	resp := decl(params, user)
	_, ok := resp.(*wishes.DeleteWishesIDOK)
	require.Equal(t, success, ok)
}

func checkSendWish(t *testing.T, user *models.UserID, wishID int64, success bool) {
	params := wishes.PutWishesIDParams{
		Content: fmt.Sprintf("wish you a %d", wishID),
		ID:      wishID,
	}
	send := api.WishesPutWishesIDHandler.Handle
	resp := send(params, user)
	_, ok := resp.(*wishes.PutWishesIDOK)
	require.Equal(t, success, ok)
}

func checkThankWish(t *testing.T, user *models.UserID, wishID int64, success bool) {
	params := wishes.PostWishesIDThankParams{ID: wishID}
	send := api.WishesPostWishesIDThankHandler.Handle
	resp := send(params, user)
	_, ok := resp.(*wishes.PostWishesIDThankNoContent)
	require.Equal(t, success, ok)
}

func TestWishes(t *testing.T) {
	req := require.New(t)

	for i := 0; i < 4; i++ {
		_, found := wishesImpl.LastCreatedWish(db, userIDs[i])
		req.False(found)
	}

	_, err := db.Exec(
		`
UPDATE users
SET karma = 1,
    created_at = NOW() - interval '4 days',
    verified = TRUE
WHERE id IN ($1, $2, $3)`,
		userIDs[0].ID, userIDs[1].ID, userIDs[2].ID)
	req.Nil(err)

	wishesImpl.CreateWishes(srv)

	id0, f0 := wishesImpl.LastCreatedWish(db, userIDs[0])
	req.True(f0)
	req.Greater(id0, int64(0))

	id1, f1 := wishesImpl.LastCreatedWish(db, userIDs[1])
	req.True(f1)
	req.Greater(id1, int64(0))

	id2, f2 := wishesImpl.LastCreatedWish(db, userIDs[2])
	req.True(f2)
	req.Greater(id2, int64(0))

	_, f3 := wishesImpl.LastCreatedWish(db, userIDs[3])
	req.False(f3)

	w0 := checkLoadWish(t, userIDs[0], id0, true)
	w1 := checkLoadWish(t, userIDs[1], id1, true)
	w2 := checkLoadWish(t, userIDs[2], id2, true)

	findReceiver := func(w *models.Wish) *models.UserID {
		req.NotNil(w.Receiver)

		for i := 0; i < 3; i++ {
			if userIDs[i].ID == w.Receiver.ID {
				return userIDs[i]
			}
		}

		return nil
	}

	r0 := findReceiver(w0)
	r1 := findReceiver(w1)
	r2 := findReceiver(w2)
	req.NotEqual(r0.ID, r1.ID)
	req.NotEqual(r1.ID, r2.ID)
	req.NotEqual(r2.ID, r0.ID)

	checkLoadWish(t, r0, id0, false)
	req.Equal(w0.State, models.WishStateNew)
	req.Equal(r0.Name, w0.Receiver.Name)

	checkDeclineWish(t, r0, id0, false)
	checkDeclineWish(t, userIDs[0], id0, true)
	checkDeclineWish(t, userIDs[0], id0, true)
	checkLoadWish(t, r0, id0, false)
	w0 = checkLoadWish(t, userIDs[0], id0, true)
	req.Equal(w0.State, models.WishStateDeclined)

	checkSendWish(t, r1, id1, false)
	checkSendWish(t, userIDs[1], id1, true)
	w1 = checkLoadWish(t, r1, id1, true)
	req.Equal(w1.State, models.WishStateSent)
	req.Nil(w1.Receiver)
	w1 = checkLoadWish(t, userIDs[1], id1, true)
	req.Equal(w1.State, models.WishStateSent)

	checkThankWish(t, userIDs[0], id0, false)
	checkThankWish(t, userIDs[1], id0, false)
	checkThankWish(t, userIDs[2], id0, false)
	checkThankWish(t, userIDs[1], id1, false)
	checkThankWish(t, r1, id1, true)
	w1 = checkLoadWish(t, r1, id1, true)
	req.Equal(w1.State, models.WishStateThanked)
	w1 = checkLoadWish(t, userIDs[1], id1, true)
	req.Equal(w1.State, models.WishStateSent)

	checkSendWish(t, userIDs[2], id2, true)

	wishesImpl.CreateWishes(srv)

	for i := 0; i < 4; i++ {
		_, found := wishesImpl.LastCreatedWish(db, userIDs[i])
		req.False(found)
	}

	_, err = db.Exec(
		`
UPDATE users
SET karma = 0,
    created_at = NOW() - interval '1 second',
    verified = FALSE
WHERE id IN ($1, $2, $3)`,
		userIDs[0].ID, userIDs[1].ID, userIDs[2].ID)
	req.Nil(err)
}
