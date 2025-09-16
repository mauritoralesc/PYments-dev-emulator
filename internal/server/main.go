package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewMainServer(port int, dashboard bool) *http.Server {
	if !gin.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas principales
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Payment Emulator Dashboard",
			"ports": []int{8001, 8002}, // Puertos de plugins
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "payment-emulator",
		})
	})

	// API para obtener estado de plugins
	r.GET("/api/plugins", func(c *gin.Context) {
		// TODO: Implementar lógica real
		c.JSON(http.StatusOK, gin.H{
			"plugins": []gin.H{
				{"name": "bancard", "port": 8001, "status": "running"},
				{"name": "pagopar", "port": 8002, "status": "running"},
			},
		})
	})

	// Cargar templates HTML embebidos
	loadTemplates(r)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
}
