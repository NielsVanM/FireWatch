package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server is an object providing easy interaction with the MUX Router
type Server struct {
	ListenPort int

	routers      []*Router
	masterRouter *mux.Router
}

// Router an representation of the mux.router with extra functionality
type Router struct {
	Name string

	endpoints  []*Endpoint
	middleware []mux.MiddlewareFunc
	router     *mux.Router
}

// Endpoint represents a URL to function map
type Endpoint struct {
	URL  string
	Func func(w http.ResponseWriter, r *http.Request)
}

// NewServer is the constructor of the Server struct
func NewServer(port int) *Server {
	s := Server{
		port,
		[]*Router{},
		mux.NewRouter(),
	}

	return &s
}

// AddRouter adds a router to the server with the prefix as url and the selected middlware
func (s *Server) AddRouter(name, prefix string) *Router {
	// Create mux subrouter linked to the master router
	mr := s.masterRouter.PathPrefix(prefix).Subrouter()

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

// AddMiddlewware adds middleware to the router
func (r *Router) AddMiddlewware(function ...mux.MiddlewareFunc) {
	r.middleware = append(r.middleware, function...)
}

// AddEndpoint adds a handlefunc to the router with the url and the function
func (r *Router) AddEndpoint(url string, function func(w http.ResponseWriter, r *http.Request)) {
	// Create new endpoint and add to router endpoints
	e := Endpoint{
		url, function,
	}

	r.endpoints = append(r.endpoints, &e)
}

// Start loads the routers and endpoints and starts the actual server
func (s *Server) Start() {
	log.Info("Starting webserver")

	// For every router and every endpoint within the router add the endpoint
	// to the router by using the handlfunc function
	for _, router := range s.routers {
		log.Info("Registering middlware for", router.Name)
		for _, middle := range router.middleware {
			router.router.Use(middle)
		}

		log.Info("Registering routes for", router.Name)
		for _, endp := range router.endpoints {
			log.Info("Registered:", endp.URL)
			router.router.HandleFunc(endp.URL, endp.Func)
		}
		fmt.Println()
	}

	s.masterRouter.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	http.ListenAndServe(fmt.Sprintf(":%d", s.ListenPort), s.masterRouter)
}
