package server

import "github.com/gin-gonic/gin"

func registerRoutes(r *gin.Engine, s *Server) {
	// v1 auth
	auth := r.Group("/v1/auth" ) 
	{
		auth.GET("/user")
		auth.POST("/user")
		auth.PUT("/user")
	}
}