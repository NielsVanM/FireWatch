package routes

import (
	"github.com/nielsvanm/firewatch/core/server"
	"github.com/nielsvanm/firewatch/views"
)

// UnprotectedRoutes is a collection of routes not protected by the authorization
// middleware
var UnprotectedRoutes = server.NewRoute("", nil,
	server.NewRoute("/auth/", nil,
		server.NewRoute("login/", views.LoginView, nil),
		server.NewRoute("logout/", views.LogoutView, nil),
	),
)

// ProtectedRoutes is a collection of routes protected by the authorization middleware
var ProtectedRoutes = server.NewRoute("", nil,
	server.NewRoute("/", views.DashboardView, nil),
	server.NewRoute("/account/", views.AccountView,
		server.NewRoute("all-device-logout/", views.LogOutAllDevicesView, nil),
		server.NewRoute("change-password/", views.ChangePasswordView, nil),
		server.NewRoute("delete/", views.DeleteAccountView, nil),
	),
	server.NewRoute("/device/", views.DeviceOverview, nil),
)
