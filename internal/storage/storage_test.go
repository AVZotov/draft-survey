package storage

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/internal/vessel"
)

const id = "123-456-789"

func getSurvey() *types.Survey {
	return &types.Survey{
		ID: "123456",
		Job: types.Job{
			JobNumber: "123456",
			DSNumber:  123456,
			Principal: "testPrincipal",
		},
		CargoOperation: types.CargoOperation{
			PlaceOfInspection: types.PortInfo{},
			Destination:       types.PortInfo{},
			Operation:         "testOperation",
			Cargo:             "testCargo",
			Packing:           "testPacking",
		},
		VesselData: vessel.VesselData{},
	}
}

func getSurveys() []*types.Survey {
	var surveys []*types.Survey
	for range 4 {
		surveys = append(surveys, getSurvey())
	}
	return surveys
}

func getUser() *types.User {
	return &types.User{
		LastName:   "Dow",
		FirstName:  "John",
		Company:    "NoName",
		Position:   "Manager",
		EmployeeID: "12345",
	}
}

func TestJSONStore_SaveAndGet(t *testing.T) {
	dir := t.TempDir()
	surveyExpected := getSurvey()
	surveyExpected.ID = id
	store := SurveyStore{Path: dir, TempPath: dir}
	if err := store.Save(surveyExpected); err != nil {
		t.Fatal(err)
	}
	surveyGot, err := store.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(surveyExpected, surveyGot) {
		t.Errorf("Expected %v, got %v", surveyExpected, surveyGot)
	}
}

func TestJSONStore_Delete(t *testing.T) {
	dir := t.TempDir()
	surveyExpected := getSurvey()
	surveyExpected.ID = id
	store := SurveyStore{Path: dir, TempPath: dir}
	if err := store.Save(surveyExpected); err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(id); err != nil {
		t.Fatal(err)
	}

	surveyGot, err := store.Get(id)
	if err == nil {
		t.Errorf("Expected error, got %#v", surveyGot)
	}
}

func TestJSONStore_GetAll(t *testing.T) {
	dir := t.TempDir()
	surveysExpected := getSurveys()
	store := SurveyStore{Path: dir, TempPath: dir}

	for i, survey := range surveysExpected {
		survey.ID = id + strconv.Itoa(i)
		if err := store.Save(survey); err != nil {
			t.Fatal(err)
		}
	}

	surveysGot, err := store.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(surveysExpected) != len(surveysGot) {
		t.Errorf("Slice length expected %d, got %d", len(surveysExpected), len(surveysGot))
	}

	for i, survey := range surveysExpected {
		if !reflect.DeepEqual(survey, surveysGot[i]) {
			t.Errorf("Expected %v, got %v", surveysExpected[i], surveysGot[i])
		}
	}
}

func TestUserStore_SaveAndGet(t *testing.T) {
	dir := t.TempDir()
	userExpected := getUser()
	store := UserStore{Path: dir}
	if err := store.Save(userExpected); err != nil {
		t.Fatal(err)
	}
	userGot, err := store.Get()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userExpected, userGot) {
		t.Errorf("Expected %v, got %v", userExpected, userGot)
	}
}

func TestUserStore_GetWithNoUser(t *testing.T) {
	dir := t.TempDir()
	store := UserStore{Path: dir}
	_, err := store.Get()
	if err == nil {
		t.Fatal("User store should not have any user")
	}
}

func TestUserStore_Delete(t *testing.T) {
	dir := t.TempDir()
	userExpected := getUser()
	store := UserStore{Path: dir}
	if err := store.Save(userExpected); err != nil {
		t.Fatal(err)
	}

	if err := store.Delete(); err != nil {
		t.Fatal(err)
	}

	userGot, err := store.Get()
	if err == nil {
		t.Errorf("Expected error, got %#v", userGot)
	}
}
