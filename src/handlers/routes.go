package handlers

import (
	"api/src/middleware"
	"api/src/notification"
	"api/src/repo"

	"github.com/gofiber/fiber/v2"
)

const ERROR = "error"
const MESSAGE = "message"

type Instance struct {
	Middleware   middleware.Instance
	Notification notification.Instance
	Repo         repos.Instance
}

func (i *Instance) Setup(app *fiber.App) {
	// go wsWorker() // run websockt worker

	app.Static("/public", "./public")

	api := app.Group("/api")
	// api.Use(func(c *fiber.Ctx) error {
	// 	return c.Status(fiber.StatusNotFound).JSON(jsonMsg("Bad request"))
	// })

	api.Post("/v1/sign-up", i.signupUser)
	api.Post("/v1/login", i.loginUser)
	api.Post("/v1/sign-up/verify", i.signupVerifyEmail)
	// api.Get("/v1/login/otp/verify", i.verifyMobileLogin)
	// api.Get("/v1/sign-up/link/mobile", i.signUpLinkMobileNumber)
	// api.Get("/v1/mobile/update/user/link", i.updateUserMobileNumber)
	// api.Get("/v1/mobile/otp/verify/link", i.verifyUserMobileLink)
	api.Get("/v1/password/reset/request", i.getPasswordResetCode)
	api.Get("/v1/password/reset/verify", i.verifyPasswordResetCode)
	api.Post("/v1/password/reset", i.resetPassword)

	api.Post("/v1/refresh-session", i.refreshSession)

	api.Get("/status", i.status)
	api.Get("/test-mail", i.testMail)

	// Auth Middleware
	api.Use(i.Middleware.AuthenticateAuthToken)

	api.Get("/v1/user", i.getUser)
	// api.Patch("/v1/user", i.updateUser)

	api.Get("/token", i.token)
}
