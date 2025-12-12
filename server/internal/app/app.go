package app

import (
	"fmt"
	"log"
	"NodeJsshell/internal/app/middleware"
	"NodeJsshell/internal/app/routes"
	"NodeJsshell/internal/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	config *config.Config
	router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	
	return &App{
		config: cfg,
		router: router,
	}
}

func (a *App) Run() error {
	routes.SetupRoutes(a.router)
	addr := fmt.Sprintf("%s:%s", a.config.Host, a.config.Port)
	log.Printf("Server starting on %s", addr)
	return a.router.Run(addr)
}

