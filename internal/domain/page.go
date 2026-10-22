package domain

import "strconv"

// Link is one navigation entry.
type Link struct {
	Href    string
	Label   string
	Current bool
}

// View is the data passed to a page template.
type View struct {
	Root        string
	Page        string
	Brand       string
	Tag         string
	Title       string
	Description string
	Nav         []Link
}

// Page is a rendered document: template name, HTTP status, and view data.
type Page struct {
	Name   string
	Status int
	View   View
}

const (
	// Brand is the placeholder site name for this template.
	Brand = "Rendering"
	// Tag is the placeholder site tagline.
	Tag = "TEMPLATE"
)

// IndexView is the default home page data.
func IndexView() View {
	return View{
		Root:        "/",
		Page:        "home",
		Brand:       Brand,
		Tag:         Tag,
		Title:       "Rendering",
		Description: "Template for a new rendering service.",
		Nav: []Link{
			{Href: "/", Label: "Home", Current: true},
		},
	}
}

// StatusView is the data for an error page.
func StatusView(code int) View {
	return View{
		Root:        "/",
		Page:        "error",
		Brand:       Brand,
		Tag:         Tag,
		Title:       strconv.Itoa(code),
		Description: statusText(code),
		Nav: []Link{
			{Href: "/", Label: "Home"},
		},
	}
}

func statusText(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	default:
		return "Error"
	}
}
