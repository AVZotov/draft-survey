package web

import (
	"github.com/AVZotov/draft-survey/internal/types"
	"github.com/AVZotov/draft-survey/web/components"
)

func TanksLayoutProps(user *types.User) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Tanks Reading",
		MetaDescription: "Vessel tanks reading",
		User:            user,
		ExtraCSS:        []string{"/static/css/tanks.css"},
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

func DraftLayoutProps(user *types.User) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Drafts Reading",
		MetaDescription: "Get vessel's draft marks",
		ExtraCSS:        []string{"/static/css/draft-readings.css"},
		ExtraJS:         []string{"/static/js/draft-readings.js"},
		User:            user,
	}
}

func DraftsPageProps(survey types.Survey) components.DraftPageProps {
	return components.DraftPageProps{
		Survey: survey,
	}
}

func SurveyLayoutProps(user *types.User) components.LayoutProps {
	return components.LayoutProps{
		Title:           "New Survey",
		MetaDescription: "General Survey Info",
		ExtraCSS:        []string{"/static/css/new-survey.css"},
		User:            user,
	}
}

func SurveyPageProps(survey types.Survey, countries []types.Country) components.SurveyPageProps {
	return components.SurveyPageProps{
		Survey:    survey,
		Countries: countries,
	}
}

// All below props in todo to split layout and needed props

func DashboardProps(user *types.User, survey *types.Survey) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Dashboard",
		MetaDescription: "Main page of Draft Survey application",
		ExtraCSS:        []string{"/static/css/dashboard.css"},
		User:            user,
		Survey:          survey,
	}
}

func ResultsPageProps(user *types.User, survey *types.Survey, results *[]types.DraftResult) components.LayoutProps {
	return components.LayoutProps{
		Title:           "Results",
		MetaDescription: "Final findings assessment",
		ExtraCSS:        []string{"/static/css/results.css"},
		ExtraJS:         []string{"/static/js/results.js"},
		User:            user,
		Survey:          survey,
	}
}

var ProfilePageProps = components.LayoutProps{
	Title:           "Surveyor Profile",
	MetaDescription: "Set up your surveyor profile",
}

var BannerFileCorrupted = components.BannerProps{
	Type:    components.Warn,
	Header:  "Profile File Corrupted",
	Message: "Your profile file could not be read. Please fill in your details again.",
}
