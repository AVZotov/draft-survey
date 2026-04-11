package handler

import (
	"net/http"
	"time"

	"github.com/AVZotov/draft-survey/internal/handler/tadaptor"
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/vessel"
	"github.com/AVZotov/draft-survey/web"
	"github.com/AVZotov/draft-survey/web/templates/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) newSurvey(c *fiber.Ctx) error {
	id := uuid.New().String()
	survey := &types.Survey{ID: id, CreatedAt: time.Now()}
	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}

	slp := web.SurveyLayoutProps(user)
	spp := web.SurveyPageProps(*survey)

	return tadaptor.Render(c, pages.NewSurvey(slp, spp))
}

func (h *Handler) saveSurvey(c *fiber.Ctx) error {
	survey, err := h.getNewSurvey(c)
	if err != nil {
		return err
	}

	if err = h.surveyRepository.Save(survey); err != nil {
		return err
	}

	switch c.FormValue("next") {
	case "dashboard":
		c.Set("HX-Redirect", "/")
	case "draft":
		c.Set("HX-Redirect", "/survey/"+survey.ID+"/draft")
	case "stay":
		c.Set("HX-Redirect", "/survey/"+survey.ID)
	default:
		c.Set("HX-Redirect", "/")
	}
	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getSurvey(c *fiber.Ctx) error {
	id := c.Params("id")
	survey, err := h.surveyRepository.Get(id)
	if err != nil {
		// TODO: Show error page
		return err
	}
	user, err := h.userRepository.Get()
	if err != nil {
		return err
	}
	slp := web.SurveyLayoutProps(user)
	spp := web.SurveyPageProps(*survey)
	return tadaptor.Render(c, pages.NewSurvey(slp, spp))
}

/*
Helper function to parse survey data from the context received.
This function not a part of API
*/
func (h *Handler) getNewSurvey(c *fiber.Ctx) (*types.Survey, error) {
	user, err := h.userRepository.Get()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	if v := c.FormValue("created_at"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			if !t.IsZero() {
				createdAt = t
			}
		}
	}

	job := types.Job{
		JobNumber: c.FormValue("job_no"),
		Principal: c.FormValue("client"),
	}
	if v, _ := parseInt(c, "ds_no"); v != nil {
		job.DSNumber = *v
	}

	cargoOperation := types.CargoOperation{
		//PlaceOfInspection: c.FormValue("port"),
		//Destination:       c.FormValue("destination"),
		Operation: c.FormValue("cargo_op"),
		Cargo:     c.FormValue("cargo"),
		Packing:   c.FormValue("packing"),
	}

	vesselData := vessel.VesselData{
		Name:             c.FormValue("vessel_name"),
		Flag:             c.FormValue("flag"),
		IMO:              c.FormValue("imo"),
		VesselType:       vessel.VesselType(c.FormValue("mmc_method")),
		CorrectionMethod: vessel.CorrectionMethod(c.FormValue("corr_method")),
	}
	if v, _ := parseInt(c, "built"); v != nil {
		vesselData.BuiltYear = *v
	}
	if v, _ := parseFloat(c, "lightship"); v != nil {
		vesselData.Lightship = *v
	}
	if v, _ := parseFloat(c, "breadth"); v != nil {
		vesselData.Breadth = *v
	}
	if v, _ := parseFloat(c, "depth"); v != nil {
		vesselData.Depth = *v
	}
	if v, _ := parseFloat(c, "lbp"); v != nil {
		vesselData.LBP = *v
	}
	if v, _ := parseFloat(c, "summer_draft"); v != nil {
		vesselData.SummerDraft = *v
	}
	if v, _ := parseFloat(c, "summer_dwt"); v != nil {
		vesselData.SummerDWT = *v
	}
	if v, _ := parseFloat(c, "summer_tpc"); v != nil {
		vesselData.SummerTPC = *v
	}
	if v, _ := parseFloat(c, "summer_freeboard"); v != nil {
		vesselData.SummerFreeboard = *v
	}
	if v, _ := parseFloat(c, "distance_pp_fwd"); v != nil {
		vesselData.DistancePPFwd = *v
	}
	if v, _ := parseFloat(c, "distance_pp_mid"); v != nil {
		vesselData.DistancePPMid = *v
	}
	if v, _ := parseFloat(c, "distance_pp_aft"); v != nil {
		vesselData.DistancePPAft = *v
	}
	if v, _ := parseFloat(c, "keel_fwd"); v != nil {
		vesselData.KeelFwd = *v
	}
	if v, _ := parseFloat(c, "keel_mid"); v != nil {
		vesselData.KeelMid = *v
	}
	if v, _ := parseFloat(c, "keel_aft"); v != nil {
		vesselData.KeelAft = *v
	}

	seaCondition := types.SeaCondition{
		Type: types.SeaConditionType(c.FormValue("sea_type")),
		Wave: types.WaveCondition(c.FormValue("sea_condition")),
		Ice:  types.IceCondition(c.FormValue("ice_condition")),
	}
	survey := &types.Survey{
		Surveyor:       *user,
		Status:         types.SurveyStatusDraft,
		ID:             c.FormValue("survey_id"),
		CreatedAt:      createdAt,
		Job:            job,
		CargoOperation: cargoOperation,
		VesselData:     vesselData,
		SeaCondition:   seaCondition,
		Remarks:        c.FormValue("remarks"),
	}

	//TODO: Validation of fields

	return survey, nil
}
