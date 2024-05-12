package handlers

import "github.com/gofiber/fiber/v2"

func (i *Instance) status(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "up"})
}

func (i *Instance) token(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"token": "valid"})
}
