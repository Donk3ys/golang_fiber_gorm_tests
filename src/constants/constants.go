package constants

import (
	"api/src/util"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

const REQ_USER_ID = "uid"
const ACCESS_TOKEN_HEADER = "Authorization"
const ACCESS_TOKEN_EXPIRY_HEADER = "Auth-Expiry"
const REFRESH_TOKEN_HEADER = "Refresh-Token"

var ACCESS_TOKEN_SECRET string
var REFRESH_TOKEN_SECRET string

var API_URL string
var WEBSITE_URL string

var ACCESS_TOKEN_DURATION = time.Minute * 20
var REFRESH_TOKEN_DURATION = time.Hour * 24 * 60

const DATE_LAYOUT = "2006-01-02"

func SetConstantsFromEnvs(envPath string) {
	atSecret := "ACCESS_TOKEN_SECRET"
	rtSecret := "REFRESH_TOKEN_SECRET"
	apiUrl := "API_URL"
	webUrl := "WEBSITE_URL"
	atDur := "ACCESS_TOKEN_DURATION"
	rtDur := "REFRESH_TOKEN_DURATION"

	if err := godotenv.Load(envPath); err != nil {
		log.Error(err)
		panic("Environment variables not set or error parsing!")
	}
	log.Infof("Setting envs from %s", envPath)

	if ACCESS_TOKEN_SECRET = os.Getenv(atSecret); ACCESS_TOKEN_SECRET == "" {
		log.Panicf("%s not set in .env", atSecret)
	}
	if REFRESH_TOKEN_SECRET = os.Getenv(rtSecret); REFRESH_TOKEN_SECRET == "" {
		log.Panicf("%s not set in .env", rtSecret)
	}

	if API_URL = os.Getenv(apiUrl); API_URL == "" {
		log.Panicf("%s not set in .env", apiUrl)
	}
	if WEBSITE_URL = os.Getenv(webUrl); WEBSITE_URL == "" {
		log.Panicf("%s not set in .env", webUrl)
	}

	if os.Getenv(atDur) != "" {
		dur, err := util.ParseDuration(os.Getenv(atDur))
		if err != nil {
			ACCESS_TOKEN_DURATION = dur
		} else {
			log.Warnf("Using default %s not set in .env %v", atDur, err)
		}
	}
	if os.Getenv(rtDur) != "" {
		dur, err := util.ParseDuration(os.Getenv(rtDur))
		if err != nil {
			REFRESH_TOKEN_DURATION = dur
		} else {
			log.Warnf("Using default %s not set in .env %v", rtDur, err)
		}
	}
}
