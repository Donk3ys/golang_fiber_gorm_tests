package handlers

import (
	"api/src/models"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func (i *Instance) getUserFromUsername(c *fiber.Ctx) error {
	// uId := (c.Locals(constants.REQ_USER_ID)).(uuid.UUID)
	username := c.Params("username")
	if username == "" || len(username) < 3 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: "No user found! [001]"})
	}

	var user models.User
	i.Repo.Db.First(&user, "username=?", username)

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: "No user found! [002]"})
	}

	return c.JSON(&user)
}

// func (i *Instance) getUserCommentsFromUsername(c *fiber.Ctx) error {
// 	uID := middleware.UserIDFromAuthToken(c)
//
// 	username := c.Params("username")
// 	if username == "" || len(username) < 3 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ERROR: "No user found!"})
// 	}
//
// 	return c.JSON(i.Repo.GetUserComments(username, uID))
// }
