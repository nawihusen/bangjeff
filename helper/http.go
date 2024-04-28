package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// HTTPSimpleResponse is function for HTTPSimpleResponse
func HTTPSimpleResponse(c *fiber.Ctx, httpStatus int) error {
	return c.Status(httpStatus).SendString(fasthttp.StatusMessage(httpStatus))
}
