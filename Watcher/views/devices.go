package views

import (
	"net/http"

	"github.com/nielsvanm/firewatch/core/models"

	"github.com/gorilla/context"
	"github.com/nielsvanm/firewatch/core/page"
)

// DeviceOverview is the dashboard page of the devices, it shows general
// information and statistics about them. As well as some notifications
func DeviceOverview(w http.ResponseWriter, r *http.Request) {
	p := page.NewPage("components/base.html", "device/overview.html")
	p.AddContext("user", context.Get(r, "user"))
	p.AddContext("device_count", models.GetAllDeviceCount())
	p.Render(w)
}
