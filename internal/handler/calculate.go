package handler

import "github.com/gofiber/fiber/v2"

func (h *Handler) calculate(c *fiber.Ctx) error {
	_, err := getProps(h, c)
	if err != nil {
		return err
	}
	return nil
}
