package routes

import (
	"os"
	"path/filepath"
	"NodeJsshell/database"
	"NodeJsshell/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := database.NewDB()
	
	shellHandler := handlers.NewShellHandler(db)
	fileHandler := handlers.NewFileHandler(db)
	proxyHandler := handlers.NewProxyHandler(db)
	payloadHandler := handlers.NewPayloadHandler()
	
	api := router.Group("/api")
	{
		api.GET("/shells", shellHandler.List)
		api.GET("/shells/:id", shellHandler.Get)
		api.POST("/shells", shellHandler.Create)
		api.PUT("/shells/:id", shellHandler.Update)
		api.DELETE("/shells/:id", shellHandler.Delete)
		api.POST("/shells/:id/test", shellHandler.Test)
		api.POST("/shells/:id/execute", shellHandler.Execute)
		api.GET("/shells/:id/info", shellHandler.GetInfo)
		
		// 文件管理路由 - 使用 /shells/:id/files 格式匹配前端
		api.GET("/shells/:id/files", fileHandler.List)
		api.GET("/shells/:id/files/read", fileHandler.Read)
		api.POST("/shells/:id/files/upload", fileHandler.Upload)
		api.GET("/shells/:id/files/download", fileHandler.Download)
		api.PUT("/shells/:id/files", fileHandler.Update)
		api.DELETE("/shells/:id/files", fileHandler.Delete)
		api.POST("/shells/:id/files/mkdir", fileHandler.Mkdir)
		
		api.GET("/proxies", proxyHandler.List)
		api.POST("/proxies", proxyHandler.Create)
		api.PUT("/proxies/:id", proxyHandler.Update)
		api.DELETE("/proxies/:id", proxyHandler.Delete)
		api.POST("/proxies/:id/test", proxyHandler.Test)
		
		api.GET("/payloads/templates", payloadHandler.GetTemplates)
		api.POST("/payloads/inject", payloadHandler.Inject)
	}
	
	distPath := "../frontend/dist"
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		distPath = "./frontend/dist"
	}
	
	absDistPath, err := filepath.Abs(distPath)
	if err != nil {
		absDistPath = distPath
	}
	
	router.Static("/assets", filepath.Join(absDistPath, "assets"))
	router.StaticFile("/favicon.ico", filepath.Join(absDistPath, "favicon.ico"))
	router.StaticFile("/vite.svg", filepath.Join(absDistPath, "vite.svg"))
	indexPath := filepath.Join(absDistPath, "index.html")
	router.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})
	router.NoRoute(func(c *gin.Context) {
		c.File(indexPath)
	})
}

