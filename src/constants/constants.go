package constants

import (
	"time"
)

const REQ_USER_ID = "req_user_id"
const AUTH_HEADER = "bearer"
const AUTH_SECRET = "BEARER_SECRET"

var API_URL string
var WEBSITE_URL string

const AUTH_TOKEN_DURATION = time.Minute * 20
const AUTH_SESSION_DURATION = time.Hour * 24 * 60

const DATE_LAYOUT = "2006-01-02"
