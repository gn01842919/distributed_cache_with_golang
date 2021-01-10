package http

import (
	"net/http"

	"../cache"
)

// Server cache
type Server struct {
	cache.Cache
}

// Listen to port
func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":12345", nil)
}

// New server
func New(c cache.Cache) *Server {
	return &Server{c} // TODO: Why returns pointer here but not in New inside cache/cache.go??
}
