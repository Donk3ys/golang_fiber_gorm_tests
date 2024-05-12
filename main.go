package main

import (
	"api/src/handlers"
	"api/src/middleware"
	"api/src/notification"
	"api/src/repo"
	"api/src/storage"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func main() {
	envPath := ".env-dev"
	if os.Getenv("BUILD") == "prod" {
		envPath = ".env-prod"
	} else if os.Getenv("BUILD") == "stage" {
		envPath = ".env-stage"
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Error(err)
		panic("Environment variables not set or error parsing!")
	}
	log.Infof("Loaded %s", envPath)

	db := storage.ConnectPostgres()
	storage.AutoMigratePostgres(db)

	storage.Seed(db)

	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Use(cors.New())

	repo := repos.Instance{
		Db: db,
		Fs: storage.FileSystem{
			Client: storage.NewLocalStorage(),
		},
	}

	nofification := notification.Instance{
		Email: notification.Email{
			Client: notification.SMTP{
				User:     os.Getenv("SMTP_USER"),
				Password: os.Getenv("SMTP_PASSWORD"),
				Host:     os.Getenv("SMTP_HOST"),
				Port:     "587",
			},
		},
		SMS: notification.SMS{
			Client: notification.Twillio{
				Username: os.Getenv("TWILIO_ACCOUNT_SID"),
				Password: os.Getenv("TWILIO_AUTH_TOKEN"),
			},
		},
	}

	appHandlers := handlers.Instance{
		Middleware:   middleware.Instance{Repo: repo},
		Notification: nofification,
		Repo:         repo,
	}
	appHandlers.Setup(app)

	c := cron.New()
	c.AddFunc("0 0 * * *", appHandlers.RemoveExpiredSessions) // 0:00 AM every day
	c.Start()

	app.Listen(":8000")
}
