package bancard

import (
	"payment-emulator/internal/plugins"
)

// BancardFactory implementa PluginFactory para crear instancias de Bancard
type BancardFactory struct{}

// NewBancardFactory crea una nueva factory para Bancard
func NewBancardFactory() plugins.PluginFactory {
	return &BancardFactory{}
}

// CreatePlugin crea una nueva instancia del plugin Bancard
func (f *BancardFactory) CreatePlugin(config *plugins.Plugin) plugins.PaymentPlugin {
	return NewBancardPlugin(config)
}

// GetPluginType devuelve el tipo de plugin que esta factory puede crear
func (f *BancardFactory) GetPluginType() string {
	return "Bancard VPOS"
}

// Función de inicialización para registrar la factory automáticamente
func init() {
	plugins.RegisterGlobalFactory("Bancard VPOS", NewBancardFactory())
	plugins.RegisterGlobalFactory("bancard", NewBancardFactory()) // Alias para compatibilidad
}
