package basicrequest

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Index Handler")
}

func PostHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Post Handler")
}

func PutHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Put Handler")
}

func DeleteHandler(c *fiber.Ctx, db *sql.DB) error {
	return c.SendString("Hello Delete Handler")
}
