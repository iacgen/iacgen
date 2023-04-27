package httpserver

import (
	"context"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port string) *HTTPServer {
	router := mux.NewRouter()
	server := &HTTPServer{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}

	router.MethodNotAllowedHandler = router.NotFoundHandler
	server.addRoutes(router)
	return server
}

func (s *HTTPServer) addRoutes(r *mux.Router) {
	h := NewAPIHandler()
	r.Methods(http.MethodGet).Path("/health").HandlerFunc(h.Health)

	apiSubRouter := r.PathPrefix("/v1/api").Subrouter()
	apiSubRouter.Use(handlers.RecoveryHandler())
	apiSubRouter.Methods(http.MethodPost).Path("/iac/generate").HandlerFunc(h.GenerateIac)
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Close() error {
	return s.server.Close()
}
