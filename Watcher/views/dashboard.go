package views

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/nielsvanm/firewatch/internal/page"
)

func DashboardView(w http.ResponseWriter, r *http.Request) {
	p := page.NewPage("components/base.html", "dashboard.html")
	p.AddContext("user", context.Get(r, "user"))
	p.AddContext("devices", "test")
	p.Render(w)
}
