package server

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, s *ServerContext) {
	userCtrl := &controllers.UserContext{
        Users:    s.Users,
        Products: s.Products,
        Orders:   s.Orders,
    }
	// v1 auth
	auth := r.Group("/v1/auth" ) 
	{
		auth.GET("/user", userCtrl.GetUser())
		auth.POST("/user", userCtrl.CreatingUser())
		auth.PUT("/user", userCtrl.UpdatingUser())
	}
}