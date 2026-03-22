package server

import (
	"backend/src/controllers"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, s *ServerContext) {
	userCtrl := &controllers.UserContext{
        Users: s.Users,
        Products: s.Products,
        Orders: s.Orders,
		Authc: s.AuthContext,
    }
	// v1 auth
	auth := r.Group("/v1/auth" ) 
	{
		auth.GET("/user", userCtrl.Login())
		auth.POST("/user", userCtrl.Register())
		auth.PUT("/user", userCtrl.Update())
	}
}