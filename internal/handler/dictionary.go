package handler

import (
	"strconv"

	"github.com/AVZotov/draft-survey/internal/constants"
	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web/components"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) seaOptions(c *fiber.Ctx) error {
	seaCtx := c.Get(constants.HXSeaType)
	if seaCtx == "" {
		return ErrUndefinedHXContext
	}

	draftIndexStr := c.Get(constants.HXDraftIndex)
	var draftIndex int
	var err error
	var scp components.SeaConditionProps

	if draftIndexStr != "" {
		draftIndex, err = strconv.Atoi(draftIndexStr)
		if err != nil {
			return err
		}
		scp.DraftIndex = &draftIndex
	}

	switch seaCtx {
	case string(types.SeaConditionTypeWave):
		scp.SeaCondition.Type = types.SeaConditionTypeWave
		return tadaptor.Render(c, components.SeaSelect(scp))
	case string(types.SeaConditionTypeIce):
		scp.SeaCondition.Type = types.SeaConditionTypeIce
		return tadaptor.Render(c, components.SeaSelect(scp))
	default:
		return ErrEmptyField
	}
}

func (h *Handler) countries(c *fiber.Ctx) error {
	countries, err := h.dictionariesRepository.GetCountries()
	if err != nil {
		return err
	}
	return tadaptor.Render(c, components.CountrySelect(components.CountrySelectProps{
		Countries: *countries,
	}))
}

func (h *Handler) ports(c *fiber.Ctx) error {
	countryCode := c.Query("country")
	if countryCode == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	ports, err := h.dictionariesRepository.GetPorts()
	if err != nil {
		return err
	}

	filtered := make([]types.Port, 0)
	for _, p := range *ports {
		if p.CountryCode == countryCode {
			filtered = append(filtered, p)
		}
	}

	return tadaptor.Render(c, components.PortSelect(components.PortSelectProps{
		Ports: filtered,
	}))
}
