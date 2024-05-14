package constants

import (
	"api/src/util"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

const REQ_USER_ID = "req_user_id"
const AUTH_HEADER = "bearer"
const AUTH_SECRET = "BEARER_SECRET"

var API_URL string
var WEBSITE_URL string

var AUTH_TOKEN_DURATION = time.Minute * 20
var AUTH_SESSION_DURATION = time.Hour * 24 * 60

const DATE_LAYOUT = "2006-01-02"

func SetConstantsFromEnvs(envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		log.Error(err)
		panic("Environment variables not set or error parsing!")
	}
	log.Infof("Setting envs from %s", envPath)

	if os.Getenv("API_URL") != "" {
		API_URL = os.Getenv("API_URL")
	} else {
		panic("API_URL not set in .env")
	}
	if os.Getenv("WEBSITE_URL") != "" {
		WEBSITE_URL = os.Getenv("WEBSITE_URL")
	} else {
		panic("WEBSITE_URL not set in .env")
	}

	if os.Getenv("AUTH_TOKEN_DURATION") != "" {
		dur, err := util.ParseDuration(os.Getenv("AUTH_TOKEN_DURATION"))
		if err != nil {
			AUTH_TOKEN_DURATION = dur
		} else {
			log.Warnf("Using default AUTH_TOKEN_DURATION not set in .env %v", err)
		}
	}
	if os.Getenv("AUTH_SESSION_DURATION") != "" {
		dur, err := util.ParseDuration(os.Getenv("AUTH_SESSION_DURATION"))
		if err != nil {
			AUTH_SESSION_DURATION = dur
		} else {
			log.Warnf("Using default AUTH_SESSION_DURATION not set in .env %v", err)
		}
	}
}
