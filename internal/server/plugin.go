package server

import (
	"fmt"
	"net/http"
	"payment-emulator/internal/plugins"

	// Importar plugins para registrar sus factories
	_ "payment-emulator/plugins/bancard"
	_ "payment-emulator/plugins/pagopar"

	"github.com/gin-gonic/gin"
)

func NewPluginServer(pluginName string, port int) *http.Server {
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

	// Cargar templates HTML embebidos PRIMERO
	loadTemplates(r)

	// Cargar y configurar el plugin específico
	if err := setupPlugin(r, pluginName, port); err != nil {
		fmt.Printf("Error setting up plugin '%s': %v\n", pluginName, err)
		// Continuar con configuración básica si hay error
		setupFallbackRoutes(r, pluginName, port)
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
}

// setupPlugin configura un plugin específico usando el registry
func setupPlugin(r *gin.Engine, pluginName string, port int) error {
	registry := plugins.GetGlobalRegistry()

	// Intentar obtener plugin del registry
	plugin, err := registry.GetPlugin(pluginName)
	if err != nil {
		// Si no existe, intentar cargarlo desde configuración
		err = registry.LoadPluginFromConfig(pluginName)
		if err != nil {
			// Si falla, crear plugin por defecto
			config := plugins.GetDefaultPlugin(pluginName, port)
			plugin, err = registry.CreatePlugin(config)
			if err != nil {
				return fmt.Errorf("failed to create plugin: %w", err)
			}
		} else {
			// Obtener el plugin recién cargado
			plugin, err = registry.GetPlugin(pluginName)
			if err != nil {
				return fmt.Errorf("failed to get loaded plugin: %w", err)
			}
		}
	}

	// Configurar rutas específicas del plugin
	plugin.SetupRoutes(r)

	// Cargar templates específicos del plugin
	loadPluginTemplates(r, plugin)

	// Ruta de documentación del plugin
	r.GET("/", func(c *gin.Context) {
		// Cargar configuración para mostrar en docs
		config, err := plugins.LoadPlugin(pluginName)
		if err != nil {
			config = plugins.GetDefaultPlugin(pluginName, port)
		}

		c.HTML(http.StatusOK, "plugin_docs.html", gin.H{
			"plugin": config,
		})
	})

	fmt.Printf("Plugin '%s' setup completed successfully\n", plugin.GetName())
	return nil
}

// loadPluginTemplates carga los templates específicos de un plugin
func loadPluginTemplates(r *gin.Engine, plugin plugins.PaymentPlugin) {
	templates := plugin.GetTemplates()

	// Por ahora, los templates están incluidos en los templates generales
	// En el futuro se pueden cargar dinámicamente aquí
	for templateName := range templates {
		fmt.Printf("Template '%s' available for plugin '%s'\n", templateName, plugin.GetName())
	}
}

// setupFallbackRoutes configura rutas básicas cuando falla la carga del plugin
func setupFallbackRoutes(r *gin.Engine, pluginName string, port int) {
	// Configuración fallback para plugins conocidos
	config := plugins.GetDefaultPlugin(pluginName, port)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "plugin_docs.html", gin.H{
			"plugin": config,
		})
	})

	// Rutas genéricas para cualquier plugin
	for _, route := range config.Routes {
		switch route.Method {
		case "POST":
			r.POST(route.Path, func(c *gin.Context) {
				handleGenericPaymentRequest(c, config, &route)
			})
		case "GET":
			r.GET(route.Path, func(c *gin.Context) {
				handleGenericPaymentRequest(c, config, &route)
			})
		}
	}

	fmt.Printf("Fallback routes configured for plugin '%s'\n", pluginName)
}

// handleGenericPaymentRequest maneja peticiones genéricas cuando no hay plugin específico
func handleGenericPaymentRequest(c *gin.Context, plugin *plugins.Plugin, route *plugins.Route) {
	// Mostrar interface de simulación según el tipo de plugin
	if plugin.Type == "iframe" {
		// Para iframe (como Bancard)
		c.HTML(http.StatusOK, "iframe_emulator.html", gin.H{
			"plugin": plugin,
			"route":  route,
			"params": c.Request.URL.Query(),
		})
	} else if plugin.Type == "popup" {
		// Para popup (como Pagopar)
		c.HTML(http.StatusOK, "popup_emulator.html", gin.H{
			"plugin": plugin,
			"route":  route,
			"params": c.Request.URL.Query(),
		})
	} else {
		// Respuesta JSON genérica
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Generic response from %s plugin", plugin.Name),
			"plugin":  plugin.Name,
			"route":   route.Path,
			"method":  route.Method,
		})
	}
}
