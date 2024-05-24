package handlers

import "github.com/gofiber/fiber/v2"

func (i *Instance) status(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "up"})
}

func (i *Instance) token(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"token": "valid"})
}

func (i *Instance) testMail(c *fiber.Ctx) error {
	err := i.Notification.Email.SendUserVerificationCode("test@email.com", 00000)
	if err != nil {
		c.JSON(fiber.Map{"error": err})
	}

	return c.JSON(fiber.Map{"status": "completed"})
}
