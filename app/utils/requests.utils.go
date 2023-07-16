package utils

import (
	fiber "github.com/gofiber/fiber/v2"
)

/*
Get request header id. From fiber context request.

Return string id
*/
func GetRequestHeaderId(c *fiber.Ctx) string {
	headers := c.GetReqHeaders()
	id := headers["Id"]
	return id
}
