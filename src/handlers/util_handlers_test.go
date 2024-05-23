package handlers_test

import (
	"api/src/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	setupTest()
	defer tearDownTest()

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	// req.Header.Set("Content-type", "application/json")
	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	json := util.JsonMapFromBytes(body)
	t.Log(json)

	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		assert.Equal(t, "up", json["status"])
		// assert.Equal(t, "alice@realworld.io", json["email"])
		// assert.Nil(t, json["bio"])
		// assert.Nil(t, json["image"])
		// assert.NotEmpty(t, json["token"])
	}
}
