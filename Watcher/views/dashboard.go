package views

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/nielsvanm/firewatch/core/page"
)

// DashboardView shows a general overview of the application status
func DashboardView(w http.ResponseWriter, r *http.Request) {
	p := page.NewPage("components/base.html", "dashboard.html")
	p.AddContext("user", context.Get(r, "user"))
	p.AddContext("devices", "test")
	p.Render(w)
}
