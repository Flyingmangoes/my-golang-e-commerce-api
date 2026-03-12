package server

import (
	"backend/config"
	"backend/middlewares"
	"backend/models/services"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Server struct {
	users 		services.UserStoreInterface
	products 	services.ProductStoreInterface 
	orders 		services.OrderStoreInterface
}

func SetupServer(us services.UserStoreInterface, prds services.ProductStoreInterface, ords services.OrderStoreInterface) *Server {
	return &Server{
		users:us,
		products: prds,
		orders: ords,
	}
}


func (s *Server) StartLoop(cfg *config.Application) {
		router := gin.Default()
		iRate := middlewares.NewIPRateLimit(rate.Limit(cfg.RateConf.RequestPerMinute), cfg.RateConf.Burst)

		router.Use(middlewares.CORSMiddleware())
		router.Use(iRate.RateLimiting())

		registerRoutes(router, s)

		addr := net.JoinHostPort(cfg.ServConf.Host, cfg.ServConf.Port)
		router.SetTrustedProxies([]string{addr})

		slog.Info("server starting", "addr", addr) 
		if err := router.Run(addr); err != nil {      
			slog.Error("server failed", "error", err)
		}
}