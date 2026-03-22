package server

import (
	"backend/src/config"
	"backend/src/middlewares"
	"backend/src/services"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Always Store the required Interface in the server context
// to make it more readable and added the new variable 
// in SetupServer

type ServerContext struct {
	Users 		services.UserStoreInterface
	Products 	services.ProductStoreInterface 
	Orders 		services.OrderStoreInterface
	AuthContext middlewares.AuthContext
}

func SetupServer(us services.UserStoreInterface, prds services.ProductStoreInterface, ords services.OrderStoreInterface, authc middlewares.AuthContext) *ServerContext {
	return &ServerContext{
		Users:us,
		Products: prds,
		Orders: ords,
		AuthContext: authc,
	}
}


func (s *ServerContext)StartLoop(cfg *config.Application) {
		router := gin.Default()
		iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RateConf.RequestPerMinute), cfg.RateConf.Burst)

		router.Use(middlewares.CORSMiddleware())
		router.Use(iRate.RateLimiting())
		router.Use(middlewares.JSONAppErrorReporter())

		registerRoutes(router, s)

		addr := net.JoinHostPort(cfg.ServConf.Host, cfg.ServConf.Port)
		router.SetTrustedProxies([]string{addr})

		slog.Info("[DEBUG] server starting", "addr", addr) 
		if err := router.Run(addr); err != nil {      
			slog.Error("server failed", "error", err)
		}
}