package web

import (
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web/components"
)

func TanksLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Tanks Reading",
		MetaDescription: "Vessel tanks reading",
		User:            user,
		ExtraCSS:        []string{"/static/css/tanks.css"},
		ExtraJS:         []string{"/static/js/tanks.js"},
		AppVersion:      version,
	}
}

func TanksPageProps(survey types.Survey, draftIndex int, trim, list *float64, trimDir, listDir string) components.TanksPageProps {
	return components.TanksPageProps{
		Survey:     survey,
		DraftIndex: draftIndex,
		Trim:       trim,
		TrimDir:    trimDir,
		List:       list,
		ListDir:    listDir,
	}
}

func DraftLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Drafts Reading",
		MetaDescription: "Get vessel's draft marks",
		ExtraCSS:        []string{"/static/css/draft-readings.css"},
		ExtraJS:         []string{"/static/js/draft-readings.js"},
		User:            user,
		AppVersion:      version,
	}
}

func DraftsPageProps(survey types.Survey) components.DraftPageProps {
	return components.DraftPageProps{
		Survey: survey,
	}
}

func SurveyLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "New Survey",
		MetaDescription: "General Survey Info",
		ExtraCSS:        []string{"/static/css/survey.css"},
		User:            user,
		AppVersion:      version,
	}
}

func SurveyPageProps(survey types.Survey) components.SurveyPageProps {
	return components.SurveyPageProps{
		Survey: survey,
	}
}

func SurveyListLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Survey List",
		MetaDescription: "List of all stored surveys",
		User:            user,
		ExtraCSS:        []string{"/static/css/survey-list.css"},
		AppVersion:      version,
	}
}

func ResultsLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Survey Results",
		MetaDescription: "Cummulative information about finilized survey",
		User:            user,
		ExtraCSS:        []string{"/static/css/results.css"},
		ExtraJS:         []string{"/static/js/results.js"},
		AppVersion:      version,
	}
}

func DashboardLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Dashboard",
		MetaDescription: "Main page of Draft Survey application",
		ExtraCSS:        []string{"/static/css/dashboard.css"},
		User:            user,
		AppVersion:      version,
	}
}

func ProfileLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Surveyor Profile",
		MetaDescription: "Surveyor profile setup page",
		ExtraCSS:        []string{"/static/css/profile.css"},
		User:            user,
		AppVersion:      version,
	}
}

func SettingsLayoutProps(user *types.User, version string) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Application Settings",
		MetaDescription: "app settings setup page",
		User:            user,
		AppVersion:      version,
	}
}
