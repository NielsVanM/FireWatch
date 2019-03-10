package server

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server is an object providing easy interaction with the MUX Router
type Server struct {
	ListenPort int

	routers      []*Router
	MasterRouter *mux.Router

	staticDirectory string
}

// Router is an implementation of the mux.router with extra functionality
type Router struct {
	Name string

	endpoints  []*Endpoint
	middleware []mux.MiddlewareFunc
	router     *mux.Router
}

// Endpoint represents a URL to function map
type Endpoint struct {
	URL    string
	Method string
	Func   func(w http.ResponseWriter, r *http.Request)
}

// Route is a representation of an endpoint, a route has SubRoutes wich append
// to the parents endpoint string. This allows for a nested overview of the
// application routes
type Route struct {
	Endpoint  string
	Method    string
	Function  func(w http.ResponseWriter, r *http.Request)
	SubRoutes []*Route
}

// NewRoute is the constructor of the type Route
func NewRoute(e, m string, f func(w http.ResponseWriter, r *http.Request), sr ...*Route) *Route {
	nr := Route{
		e, m, f, sr,
	}

	return &nr
}

// NewServer is the constructor of the Server struct
func NewServer(port int) *Server {
	s := Server{
		port,
		[]*Router{},
		mux.NewRouter(),
		"./static/",
	}

	return &s
}

// AddRouter adds a router to the server with the prefix as url and the selected middlware
func (s *Server) AddRouter(name, prefix string) *Router {
	// Create mux subrouter linked to the master router
	mr := s.MasterRouter.PathPrefix(prefix).Subrouter()

	// Create router obj
	r := Router{
		name,
		[]*Endpoint{},
		[]mux.MiddlewareFunc{},
		mr,
	}

	// Add to list
	s.routers = append(s.routers, &r)

	return &r
}

// AddMiddleware adds middleware to the router
func (r *Router) AddMiddleware(function ...mux.MiddlewareFunc) {
	r.middleware = append(r.middleware, function...)
}

// AddEndpoint adds a handlefunc to the router with the url and the function
func (r *Router) AddEndpoint(url, method string, function func(w http.ResponseWriter, r *http.Request)) {
	// Create new endpoint and add to router endpoints
	e := Endpoint{
		url, method, function,
	}

	r.endpoints = append(r.endpoints, &e)
}

// ParseRouteMap takes a route object and parses it into all the defined
// endpoints
func (r *Router) ParseRouteMap(route *Route) {
	r.parseRoute("", route)
}

func (r *Router) parseRoute(prefix string, route *Route) {
	// If the subroutes aren't nil loop over them
	for _, subRoute := range route.SubRoutes {
		if subRoute == nil {
			continue
		}
		// Check if the subRoute has a function
		if subRoute.Function != nil {
			// Create endpoint
			r.AddEndpoint(prefix+subRoute.Endpoint, subRoute.Method, subRoute.Function)
		}

		// Parse the new route
		r.parseRoute(prefix+subRoute.Endpoint, subRoute)
	}

}

// SetStaticDir sets the static file directory to a path
func (s *Server) SetStaticDir(path string) {
	s.staticDirectory = path
}

// Start loads the routers and endpoints and starts the actual server
func (s *Server) Start() {
	log.Info("Starting webserver")

	// For every router and every endpoint within the router add the endpoint
	// to the router by using the handlfunc function
	for _, router := range s.routers {
		log.Info("Registering middleware for", router.Name)
		for _, middle := range router.middleware {
			router.router.Use(middle)
		}

		log.Info("Registering routes for", router.Name)
		for _, endp := range router.endpoints {
			log.Info("Registered:", endp.URL)
			router.router.HandleFunc(endp.URL, endp.Func).
				Methods(strings.Split(endp.Method, ",")...)
		}
		fmt.Println()
	}

	s.MasterRouter.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDirectory))))

	log.Info("Starting server")
	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.ListenPort), s.MasterRouter))
}
