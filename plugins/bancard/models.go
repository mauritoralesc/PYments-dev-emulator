package bancard

// BancardOrderRequest representa una petición de pago de Bancard
type BancardOrderRequest struct {
	PublicKey string           `json:"public_key" binding:"required"`
	Operation BancardOperation `json:"operation" binding:"required"`
	ReturnURL string           `json:"return_url,omitempty"`
	CancelURL string           `json:"cancel_url,omitempty"`
}

// BancardOperation representa los datos de la operación
type BancardOperation struct {
	Token            string                 `json:"token" binding:"required"`
	ShopProcessID    string                 `json:"shop_process_id" binding:"required"`
	Amount           string                 `json:"amount" binding:"required"`
	Currency         string                 `json:"currency,omitempty"`
	NumberOfPayments int                    `json:"number_of_payments,omitempty"`
	AdditionalData   map[string]interface{} `json:"additional_data,omitempty"`
}

// BancardOrderResponse representa la respuesta de creación de orden
type BancardOrderResponse struct {
	Status      string `json:"status"`
	ProcessID   string `json:"process_id"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Message     string `json:"message,omitempty"`
}

// BancardConfirmationRequest representa una petición de confirmación
type BancardConfirmationRequest struct {
	ShopProcessID string                  `json:"shop_process_id" binding:"required"`
	Token         string                  `json:"token" binding:"required"`
	Operation     BancardConfirmOperation `json:"operation" binding:"required"`
}

// BancardConfirmOperation representa los datos de confirmación
type BancardConfirmOperation struct {
	Token string `json:"token" binding:"required"`
}

// BancardConfirmationResponse representa la respuesta de confirmación
type BancardConfirmationResponse struct {
	Status              string              `json:"status"`
	Message             string              `json:"message"`
	TransactionID       string              `json:"transaction_id,omitempty"`
	AuthorizationNumber string              `json:"authorization_number,omitempty"`
	TicketNumber        string              `json:"ticket_number,omitempty"`
	ResponseCode        string              `json:"response_code,omitempty"`
	ResponseDescription string              `json:"response_description,omitempty"`
	Amount              string              `json:"amount,omitempty"`
	Currency            string              `json:"currency,omitempty"`
	Security            BancardSecurityInfo `json:"security,omitempty"`
}

// BancardSecurityInfo representa información de seguridad
type BancardSecurityInfo struct {
	Customer     BancardCustomerInfo `json:"customer"`
	CardInfo     BancardCardInfo     `json:"card_info"`
	RiskAnalysis BancardRiskAnalysis `json:"risk_analysis"`
}

// BancardCustomerInfo representa información del cliente
type BancardCustomerInfo struct {
	Document     string `json:"document"`
	DocumentType string `json:"document_type"`
	Email        string `json:"email"`
	CellPhone    string `json:"cell_phone"`
}

// BancardCardInfo representa información de la tarjeta
type BancardCardInfo struct {
	Bin           string `json:"bin"`
	Last4         string `json:"last_4"`
	Brand         string `json:"brand"`
	Type          string `json:"type"`
	Issuer        string `json:"issuer"`
	IssuerCountry string `json:"issuer_country"`
}

// BancardRiskAnalysis representa análisis de riesgo
type BancardRiskAnalysis struct {
	Score          int    `json:"score"`
	Recommendation string `json:"recommendation"`
}

// BancardCheckoutData representa los datos para el checkout
type BancardCheckoutData struct {
	ProcessID     string              `json:"process_id"`
	Amount        string              `json:"amount"`
	Currency      string              `json:"currency"`
	ShopProcessID string              `json:"shop_process_id"`
	ReturnURL     string              `json:"return_url"`
	CancelURL     string              `json:"cancel_url"`
	OrderDetails  BancardOrderDetails `json:"order_details"`
}

// BancardOrderDetails representa los detalles de la orden
type BancardOrderDetails struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

// BancardSimulationResult representa el resultado de una simulación
type BancardSimulationResult struct {
	Status        string `json:"status"`
	ProcessID     string `json:"process_id"`
	Message       string `json:"message"`
	RedirectURL   string `json:"redirect_url,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// Constantes para Bancard
const (
	// Estados de transacción
	StatusSuccess   = "success"
	StatusError     = "error"
	StatusPending   = "pending"
	StatusCancelled = "cancelled"

	// Monedas soportadas
	CurrencyPYG = "PYG"
	CurrencyUSD = "USD"

	// Tipos de tarjeta
	CardTypeCredit = "credit"
	CardTypeDebit  = "debit"

	// Marcas de tarjeta
	BrandVisa       = "VISA"
	BrandMastercard = "MASTERCARD"
	BrandAmex       = "AMEX"

	// URLs por defecto
	DefaultReturnURL = "/bancard/return"
	DefaultCancelURL = "/bancard/cancel"
)
