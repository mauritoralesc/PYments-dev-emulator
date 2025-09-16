package pagopar

import (
	"payment-emulator/internal/plugins"
)

// PagoparFactory implementa PluginFactory para crear instancias de Pagopar
type PagoparFactory struct{}

// NewPagoparFactory crea una nueva factory para Pagopar
func NewPagoparFactory() plugins.PluginFactory {
	return &PagoparFactory{}
}

// CreatePlugin crea una nueva instancia del plugin Pagopar
func (f *PagoparFactory) CreatePlugin(config *plugins.Plugin) plugins.PaymentPlugin {
	return NewPagoparPlugin(config)
}

// GetPluginType devuelve el tipo de plugin que esta factory puede crear
func (f *PagoparFactory) GetPluginType() string {
	return "Pagopar"
}

// Función de inicialización para registrar la factory automáticamente
func init() {
	plugins.RegisterGlobalFactory("Pagopar", NewPagoparFactory())
}
