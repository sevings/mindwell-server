package test

import (
	"testing"

	"github.com/sevings/mindwell-server/models"

	"github.com/sevings/mindwell-server/restapi/operations/adm"
	"github.com/stretchr/testify/require"
)

const admRegFinished = true

func checkGrandsonAddress(t *testing.T, userID *models.UserID, name, postcode, country, address, phone, comment string, anonymous bool) {
	if admRegFinished {
		return
	}

	req := require.New(t)

	load := api.AdmGetAdmGrandsonHandler.Handle
	resp := load(adm.GetAdmGrandsonParams{}, userID)
	body, ok := resp.(*adm.GetAdmGrandsonOK)
	req.True(ok)

	addr := body.Payload
	req.Equal(name, addr.Name)
	req.Equal(postcode, addr.Postcode)
	req.Equal(country, addr.Country)
	req.Equal(address, addr.Address)
	req.Equal(phone, addr.Phone)
	req.Equal(comment, addr.Comment)
	req.Equal(anonymous, addr.Anonymous)
}

func updateGrandsonAddress(t *testing.T, userID *models.UserID, name, postcode, country, address, phone, comment string, anonymous bool) {
	if admRegFinished {
		return
	}

	params := adm.PostAdmGrandsonParams{
		Name:      name,
		Postcode:  postcode,
		Country:   country,
		Address:   address,
		Phone:     &phone,
		Comment:   &comment,
		Anonymous: &anonymous,
	}

	post := api.AdmPostAdmGrandsonHandler.Handle
	resp := post(params, userID)
	_, ok := resp.(*adm.PostAdmGrandsonOK)
	require.True(t, ok)

	checkGrandsonAddress(t, userID, name, postcode, country, address, phone, comment, anonymous)
}

func checkAdmStat(t *testing.T, grandsons int64) {
	if admRegFinished {
		return
	}

	req := require.New(t)

	load := api.AdmGetAdmStatHandler.Handle
	resp := load(adm.GetAdmStatParams{}, userIDs[0])
	body, ok := resp.(*adm.GetAdmStatOK)
	req.True(ok)

	stat := body.Payload
	req.Equal(grandsons, stat.Grandsons)
}

func TestAdm(t *testing.T) {
	checkAdmStat(t, 0)

	checkGrandsonAddress(t, userIDs[0], "", "", "", "", "", "", false)
	checkAdmStat(t, 0)

	updateGrandsonAddress(t, userIDs[0], "aaa", "213", "Aaa", "aaaa", "aaaaa", "aaaaaa", false)
	checkAdmStat(t, 1)

	updateGrandsonAddress(t, userIDs[0], "bbb", "5654", "Bbb", "bbbb", "bbbbb", "bbbbbb", true)
	checkAdmStat(t, 1)

	updateGrandsonAddress(t, userIDs[1], "vvv", "5654", "Bbb", "bbbb", "ccccc", "cccccc", true)
	checkAdmStat(t, 2)
}
