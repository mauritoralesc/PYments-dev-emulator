package plugins

import (
	"fmt"
	"sync"
)

// PluginRegistry maneja el registro y descubrimiento de plugins
type PluginRegistry struct {
	plugins   map[string]PaymentPlugin
	factories map[string]PluginFactory
	mutex     sync.RWMutex
}

// NewPluginRegistry crea una nueva instancia del registro de plugins
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins:   make(map[string]PaymentPlugin),
		factories: make(map[string]PluginFactory),
	}
}

// RegisterFactory registra una factory para crear plugins de un tipo específico
func (r *PluginRegistry) RegisterFactory(pluginType string, factory PluginFactory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.factories[pluginType] = factory
}

// RegisterPlugin registra una instancia de plugin directamente
func (r *PluginRegistry) RegisterPlugin(name string, plugin PaymentPlugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.plugins[name] = plugin
}

// GetPlugin obtiene un plugin por nombre
func (r *PluginRegistry) GetPlugin(name string) (PaymentPlugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin, nil
}

// CreatePlugin crea un plugin usando la factory apropiada
func (r *PluginRegistry) CreatePlugin(config *Plugin) (PaymentPlugin, error) {
	r.mutex.RLock()
	factory, exists := r.factories[config.Name]
	r.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no factory registered for plugin type '%s'", config.Name)
	}

	plugin := factory.CreatePlugin(config)

	// Registrar automáticamente el plugin creado
	r.RegisterPlugin(config.Name, plugin)

	return plugin, nil
}

// ListPlugins devuelve una lista de todos los plugins registrados
func (r *PluginRegistry) ListPlugins() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	return names
}

// HasPlugin verifica si un plugin está registrado
func (r *PluginRegistry) HasPlugin(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.plugins[name]
	return exists
}

// LoadPluginFromConfig carga un plugin desde su configuración YAML y lo registra
func (r *PluginRegistry) LoadPluginFromConfig(pluginName string) error {
	// Cargar configuración del plugin
	config, err := LoadPlugin(pluginName)
	if err != nil {
		// Si no se encuentra configuración, crear plugin por defecto
		config = GetDefaultPlugin(pluginName, 0) // El puerto será asignado por el servidor
	}

	// Crear plugin usando factory
	plugin, err := r.CreatePlugin(config)
	if err != nil {
		return fmt.Errorf("failed to create plugin '%s': %w", pluginName, err)
	}

	// Registrar también con el nombre solicitado (para compatibilidad)
	if pluginName != config.Name {
		r.RegisterPlugin(pluginName, plugin)
	}

	fmt.Printf("Plugin '%s' loaded successfully\n", plugin.GetName())
	return nil
}

// Global registry instance
var globalRegistry = NewPluginRegistry()

// GetGlobalRegistry devuelve la instancia global del registro
func GetGlobalRegistry() *PluginRegistry {
	return globalRegistry
}

// RegisterGlobalPlugin registra un plugin en el registro global
func RegisterGlobalPlugin(name string, plugin PaymentPlugin) {
	globalRegistry.RegisterPlugin(name, plugin)
}

// RegisterGlobalFactory registra una factory en el registro global
func RegisterGlobalFactory(pluginType string, factory PluginFactory) {
	globalRegistry.RegisterFactory(pluginType, factory)
}

// GetGlobalPlugin obtiene un plugin del registro global
func GetGlobalPlugin(name string) (PaymentPlugin, error) {
	return globalRegistry.GetPlugin(name)
}
