package server

import (
	"github.com/gin-gonic/gin"
)

// RouteFunc is a function that registers routes on a RouterGroup.
type RouteFunc func(r *gin.RouterGroup)

// SetupRoutes creates a versioned API group and registers the provided route
// functions on it. When version is non-empty the group path is "/api/<version>",
// otherwise "/api".
func (s *Server) SetupRoutes(version string, routes ...RouteFunc) {
	groupPath := "/api"
	if version != "" {
		groupPath = "/api/" + version
	}
	api := s.engine.Group(groupPath)
	for _, route := range routes {
		route(api)
	}
}

// RegisterRoutes is a convenience wrapper around SetupRoutes with an empty
// version string, resulting in routes mounted under "/api".
func (s *Server) RegisterRoutes(routes ...RouteFunc) {
	s.SetupRoutes("", routes...)
}
