package pagopar

import "time"

// PagoparOrderRequest representa la petición para iniciar transacción en Pagopar
type PagoparOrderRequest struct {
	Token          string                 `json:"token" binding:"required"`
	PublicKey      string                 `json:"public_key" binding:"required"`
	MontoTotal     string                 `json:"monto_total" binding:"required"`
	Comprador      PagoparComprador       `json:"comprador" binding:"required"`
	ComprasItems   []PagoparComprasItem   `json:"compras_items" binding:"required"`
	UrlResultado   string                 `json:"url_resultado,omitempty"`
	UrlCancelacion string                 `json:"url_cancelacion,omitempty"`
	UrlRespuesta   string                 `json:"url_respuesta,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// PagoparComprador representa los datos del comprador
type PagoparComprador struct {
	Email     string `json:"email" binding:"required"`
	Telefono  string `json:"telefono,omitempty"`
	Documento string `json:"documento,omitempty"`
	Nombre    string `json:"nombre,omitempty"`
}

// PagoparComprasItem representa un item de compra
type PagoparComprasItem struct {
	Nombre      string `json:"nombre" binding:"required"`
	Precio      string `json:"precio" binding:"required"`
	Cantidad    int    `json:"cantidad" binding:"required"`
	Descripcion string `json:"descripcion,omitempty"`
}

// PagoparOrderResponse representa la respuesta de creación de orden
type PagoparOrderResponse struct {
	Respuesta bool                 `json:"respuesta"`
	Resultado []PagoparOrderResult `json:"resultado,omitempty"`
	Error     string               `json:"error,omitempty"`
}

// PagoparOrderResult representa el resultado de la orden
type PagoparOrderResult struct {
	Data   string `json:"data"`   // Hash de la orden
	Pedido string `json:"pedido"` // Número de pedido
}

// PagoparPaymentMethodsRequest representa la petición para obtener métodos de pago
type PagoparPaymentMethodsRequest struct {
	Token        string `json:"token" binding:"required"`
	TokenPublico string `json:"token_publico" binding:"required"`
}

// PagoparPaymentMethodsResponse representa la respuesta de métodos de pago
type PagoparPaymentMethodsResponse struct {
	Respuesta bool                   `json:"respuesta"`
	Resultado []PagoparPaymentMethod `json:"resultado,omitempty"`
}

// PagoparPaymentMethod representa un método de pago disponible
type PagoparPaymentMethod struct {
	FormaPago            string `json:"forma_pago"`
	Titulo               string `json:"titulo"`
	Descripcion          string `json:"descripcion"`
	MontoMinimo          string `json:"monto_minimo"`
	PorcentajeComision   string `json:"porcentaje_comision"`
	PagosInternacionales bool   `json:"pagos_internacionales,omitempty"`
}

// PagoparOrderStatusRequest representa la petición para consultar estado de pedido
type PagoparOrderStatusRequest struct {
	HashPedido   string `json:"hash_pedido" binding:"required"`
	Token        string `json:"token" binding:"required"`
	TokenPublico string `json:"token_publico" binding:"required"`
}

// PagoparOrderStatusResponse representa la respuesta del estado de pedido
type PagoparOrderStatusResponse struct {
	Respuesta bool                     `json:"respuesta"`
	Resultado []PagoparOrderStatusData `json:"resultado,omitempty"`
}

// PagoparOrderStatusData representa los datos del estado de un pedido
type PagoparOrderStatusData struct {
	Pagado                   bool        `json:"pagado"`
	NumeroComprobanteInterno string      `json:"numero_comprobante_interno"`
	UltimoMensajeError       interface{} `json:"ultimo_mensaje_error"`
	FormaPago                string      `json:"forma_pago"`
	FechaPago                interface{} `json:"fecha_pago"`
	Monto                    string      `json:"monto"`
	FechaMaximaPago          string      `json:"fecha_maxima_pago"`
	HashPedido               string      `json:"hash_pedido"`
	NumeroPedido             string      `json:"numero_pedido"`
	Cancelado                bool        `json:"cancelado"`
	FormaPagoIdentificador   string      `json:"forma_pago_identificador"`
	Token                    string      `json:"token,omitempty"`
	MensajeResultadoPago     interface{} `json:"mensaje_resultado_pago,omitempty"`
}

// PagoparWebhookData representa los datos del webhook de Pagopar
type PagoparWebhookData struct {
	Resultado []PagoparOrderStatusData `json:"resultado"`
	Respuesta bool                     `json:"respuesta"`
}

// PagoparCheckoutData representa los datos necesarios para el checkout
type PagoparCheckoutData struct {
	Hash         string                 `json:"hash"`
	FormaPago    string                 `json:"forma_pago"`
	Methods      []PagoparPaymentMethod `json:"methods"`
	OrderDetails PagoparOrderDetails    `json:"order_details"`
}

// PagoparOrderDetails representa los detalles de la orden para mostrar en checkout
type PagoparOrderDetails struct {
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// PagoparSimulationResult representa el resultado de una simulación
type PagoparSimulationResult struct {
	Status      string      `json:"status"` // success, error, pending, cancel
	OrderHash   string      `json:"order_hash"`
	Message     string      `json:"message"`
	RedirectURL string      `json:"redirect_url,omitempty"`
	WebhookData interface{} `json:"webhook_data,omitempty"`
}

// Constantes para Pagopar
const (
	// Estados de pago
	PaymentStatusSuccess = "success"
	PaymentStatusError   = "error"
	PaymentStatusPending = "pending"
	PaymentStatusCancel  = "cancel"

	// Métodos de pago comunes
	PaymentMethodCredit    = "9"
	PaymentMethodTigoMoney = "10"
	PaymentMethodBank      = "11"
	PaymentMethodPIX       = "25"
	PaymentMethodQR        = "24"

	// Configuración por defecto
	DefaultCurrency = "PYG"
	DefaultAmount   = "100000.00"
)
