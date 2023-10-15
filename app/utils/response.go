package utils

import "github.com/gofiber/fiber/v2"

type Messages = []any

type Response struct {
	Code     int      `json:"code"`
	Messages Messages `json:"messages,omitempty"`
	Data     any      `json:"data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, resp Response) error {
	if resp.Code == 0 {
		resp.Code = fiber.StatusOK
	}
	c.Status(resp.Code)
	return c.JSON(resp)
}
