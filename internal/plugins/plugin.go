package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Plugin struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Port        int     `yaml:"port"`
	Type        string  `yaml:"type"` // "iframe" o "redirección"
	Enabled     bool    `yaml:"enabled"`
	Routes      []Route `yaml:"routes"`
}

type Route struct {
	Path         string `yaml:"path"`
	Method       string `yaml:"method"`
	ResponseType string `yaml:"response_type"`
}

func LoadPlugin(name string) (*Plugin, error) {
	configPath := filepath.Join("plugins", name, "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var plugin Plugin
	err = yaml.Unmarshal(data, &plugin)
	return &plugin, err
}

func GetDefaultPlugin(name string, port int) *Plugin {
	switch name {
	case "bancard":
		return &Plugin{
			Name:        "Bancard VPOS",
			Description: "Emulador de Bancard VPOS",
			Port:        port,
			Type:        "iframe",
			Enabled:     true,
			Routes: []Route{
				{Path: "/vpos/api/0.3/single_buy", Method: "POST", ResponseType: "redirect"},
				{Path: "/vpos/api/0.3/confirmation", Method: "POST", ResponseType: "json"},
			},
		}
	case "pagopar":
		return &Plugin{
			Name:        "Pagopar",
			Description: "Emulador de Pagopar",
			Port:        port,
			Type:        "redirect",
			Enabled:     true,
			Routes: []Route{
				{Path: "/api/comercios/2.0/iniciar-transaccion", Method: "POST", ResponseType: "json"},
				{Path: "/api/forma-pago/1.1/traer", Method: "POST", ResponseType: "json"},
				{Path: "/api/pedidos/1.1/traer", Method: "POST", ResponseType: "json"},
				{Path: "/pagos/:hash", Method: "GET", ResponseType: "html"},
			},
		}
	default:
		return &Plugin{
			Name:        name,
			Description: fmt.Sprintf("Plugin personalizado %s", name),
			Port:        port,
			Type:        "iframe",
			Enabled:     true,
			Routes: []Route{
				{Path: "/pay", Method: "POST", ResponseType: "redirect"},
			},
		}
	}
}

func GetAvailablePlugins() []Plugin {
	plugins := []Plugin{}

	// Plugins embebidos
	plugins = append(plugins, *GetDefaultPlugin("bancard", 8001))
	plugins = append(plugins, *GetDefaultPlugin("pagopar", 8002))

	// Buscar plugins personalizados
	pluginsDir := "plugins"
	if entries, err := os.ReadDir(pluginsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				if plugin, err := LoadPlugin(entry.Name()); err == nil {
					plugins = append(plugins, *plugin)
				}
			}
		}
	}

	return plugins
}

func CreatePluginTemplate(name string) error {
	pluginDir := filepath.Join("plugins", name)

	// Crear directorio
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return err
	}

	// Crear config.yaml template
	configContent := fmt.Sprintf(`name: "%s"
description: "Plugin personalizado para %s"
port: 8003
type: "iframe"  # o "popup"
enabled: true
routes:
  - path: "/pay"
    method: "POST"
    response_type: "redirect"
  - path: "/confirmation"
    method: "POST"
    response_type: "json"
`, name, name)

	configPath := filepath.Join(pluginDir, "config.yaml")
	return os.WriteFile(configPath, []byte(configContent), 0644)
}
