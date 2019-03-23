package routes

import (
	"github.com/nielsvanm/firewatch/api"
	"github.com/nielsvanm/firewatch/core/server"
)

// UnprotectedRoutes is a collection of routes not protected by the authorization
// middleware
var UnprotectedRoutes = server.NewRoute("", "GET", nil,
	server.NewRoute("/api/v1/", "POST", nil,
		server.NewRoute("login/", "POST", api.Login, nil),
		server.NewRoute("logout/", "POST", api.Logout, nil),
		server.NewRoute("verify/", "POST", api.VerifyToken, nil),
	),
)

// ProtectedRoutes is a collection of routes protected by the authorization middleware
var ProtectedRoutes = server.NewRoute("", "GET", nil,
	server.NewRoute("/api/v1/", "GET", nil,
		server.NewRoute("device/", "GET", api.GetAllDevices,
			server.NewRoute("{id}/", "GET", api.GetDevice, nil),
			server.NewRoute("{id}/", "POST", api.CreateDevice, nil),
			server.NewRoute("{id}/", "PUT", api.UpdateDevice, nil),
			server.NewRoute("{id}/", "DELETE", api.DeleteDevice, nil),
		),
		server.NewRoute("data/", "GET", api.GetAllData,
			server.NewRoute("{id}/", "GET", api.GetData, nil),
			server.NewRoute("{id}/", "POST", api.CreateData, nil),
			server.NewRoute("{id}/", "PUT", api.UpdateData, nil),
			server.NewRoute("{id}/", "DELETE", api.DeleteData, nil),
		),
	),
)
