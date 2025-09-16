package plugins

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

// PaymentPlugin define la interfaz principal para todos los plugins de pago
type PaymentPlugin interface {
	// GetName devuelve el nombre del plugin
	GetName() string

	// GetType devuelve el tipo de plugin (iframe, popup, etc.)
	GetType() string

	// SetupRoutes configura las rutas específicas del plugin
	SetupRoutes(router *gin.Engine)

	// GetTemplates devuelve los templates específicos del plugin
	GetTemplates() map[string]string

	// HandlePaymentRequest maneja las peticiones de pago genéricas
	HandlePaymentRequest(c *gin.Context, route *Route)
}

// TemplateProvider define la interfaz para proveedores de templates
type TemplateProvider interface {
	// LoadTemplates carga y devuelve un mapa de templates
	LoadTemplates() map[string]*template.Template

	// GetTemplate devuelve un template específico por nombre
	GetTemplate(name string) (*template.Template, error)
}

// PluginFactory define la interfaz para crear instancias de plugins
type PluginFactory interface {
	// CreatePlugin crea una nueva instancia del plugin
	CreatePlugin(config *Plugin) PaymentPlugin

	// GetPluginType devuelve el tipo de plugin que esta factory puede crear
	GetPluginType() string
}

// WebhookHandler define la interfaz para manejar webhooks
type WebhookHandler interface {
	// HandleWebhook procesa un webhook recibido
	HandleWebhook(c *gin.Context) error

	// ValidateWebhook valida la autenticidad del webhook
	ValidateWebhook(payload []byte, signature string) error
}

// PaymentStatus representa el estado de un pago
type PaymentStatus struct {
	Paid        bool        `json:"paid"`
	Amount      string      `json:"amount"`
	Currency    string      `json:"currency"`
	Status      string      `json:"status"`
	OrderHash   string      `json:"order_hash"`
	OrderNumber string      `json:"order_number"`
	PaymentDate interface{} `json:"payment_date"`
	Error       string      `json:"error,omitempty"`
}

// PaymentRequest representa una petición de pago genérica
type PaymentRequest struct {
	Token       string                 `json:"token"`
	PublicKey   string                 `json:"public_key"`
	Amount      string                 `json:"amount"`
	Currency    string                 `json:"currency"`
	Description string                 `json:"description"`
	OrderID     string                 `json:"order_id"`
	ReturnURL   string                 `json:"return_url"`
	CancelURL   string                 `json:"cancel_url"`
	WebhookURL  string                 `json:"webhook_url"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
