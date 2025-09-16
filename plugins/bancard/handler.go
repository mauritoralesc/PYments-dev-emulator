package bancard

import (
	"fmt"
	"math/rand"
	"net/http"
	"payment-emulator/internal/plugins"

	"github.com/gin-gonic/gin"
)

// BancardPlugin implementa PaymentPlugin para Bancard
type BancardPlugin struct {
	name       string
	pluginType string
	config     *plugins.Plugin
}

// NewBancardPlugin crea una nueva instancia del plugin de Bancard
func NewBancardPlugin(config *plugins.Plugin) plugins.PaymentPlugin {
	return &BancardPlugin{
		name:       "Bancard VPOS",
		pluginType: "iframe",
		config:     config,
	}
}

// GetName devuelve el nombre del plugin
func (p *BancardPlugin) GetName() string {
	return p.name
}

// GetType devuelve el tipo del plugin
func (p *BancardPlugin) GetType() string {
	return p.pluginType
}

// SetupRoutes configura todas las rutas específicas de Bancard
func (p *BancardPlugin) SetupRoutes(r *gin.Engine) {
	p.setupAPIRoutes(r)
	p.setupCheckoutRoutes(r)
	p.setupEmulatorRoutes(r)
}

// GetTemplates devuelve los templates específicos de Bancard
func (p *BancardPlugin) GetTemplates() map[string]string {
	return GetBancardTemplates()
}

// HandlePaymentRequest maneja peticiones de pago genéricas
func (p *BancardPlugin) HandlePaymentRequest(c *gin.Context, route *plugins.Route) {
	c.HTML(http.StatusOK, "bancard_docs.html", gin.H{
		"plugin": p.config,
		"route":  route,
		"params": c.Request.URL.Query(),
	})
}

// setupAPIRoutes configura las rutas de la API de Bancard
func (p *BancardPlugin) setupAPIRoutes(r *gin.Engine) {
	// API versión 0.3
	v03 := r.Group("/vpos/api/0.3")
	{
		v03.POST("/single_buy", p.handleSingleBuy)
		v03.POST("/confirmation", p.handleConfirmation)
		v03.GET("/single_buy/:process_id", p.handleGetTransaction)
	}

	// API legacy
	r.POST("/bancard/single_buy", p.handleSingleBuy)
	r.POST("/bancard/confirmation", p.handleConfirmation)
}

// setupCheckoutRoutes configura las rutas del checkout
func (p *BancardPlugin) setupCheckoutRoutes(r *gin.Engine) {
	r.GET("/bancard/checkout/:process_id", p.handleCheckout)
	r.GET("/bancard/return", p.handleReturn)
	r.GET("/bancard/cancel", p.handleCancel)
}

// setupEmulatorRoutes configura las rutas del emulador
func (p *BancardPlugin) setupEmulatorRoutes(r *gin.Engine) {
	r.POST("/emulator/bancard/:process_id", p.handleEmulatorPayment)
	r.GET("/emulator/bancard/result", p.handleEmulatorResult)
}

// handleSingleBuy maneja la creación de una transacción de compra simple
func (p *BancardPlugin) handleSingleBuy(c *gin.Context) {
	var request BancardOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, BancardOrderResponse{
			Status:  StatusError,
			Message: "Datos JSON inválidos: " + err.Error(),
		})
		return
	}

	// Validar campos obligatorios
	if request.PublicKey == "" {
		c.JSON(http.StatusBadRequest, BancardOrderResponse{
			Status:  StatusError,
			Message: "public_key es obligatorio",
		})
		return
	}

	if request.Operation.Token == "" {
		c.JSON(http.StatusBadRequest, BancardOrderResponse{
			Status:  StatusError,
			Message: "operation.token es obligatorio",
		})
		return
	}

	if request.Operation.Amount == "" {
		c.JSON(http.StatusBadRequest, BancardOrderResponse{
			Status:  StatusError,
			Message: "operation.amount es obligatorio",
		})
		return
	}

	// Generar ProcessID único
	processID := generateProcessID()
	redirectURL := fmt.Sprintf("/bancard/checkout/%s", processID)

	response := BancardOrderResponse{
		Status:      StatusSuccess,
		ProcessID:   processID,
		RedirectURL: redirectURL,
		Message:     "Transacción creada exitosamente",
	}

	c.JSON(http.StatusOK, response)
}

// handleConfirmation maneja la confirmación de una transacción
func (p *BancardPlugin) handleConfirmation(c *gin.Context) {
	var request BancardConfirmationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, BancardConfirmationResponse{
			Status:  StatusError,
			Message: "Datos JSON inválidos: " + err.Error(),
		})
		return
	}

	// Validar campos obligatorios
	if request.ShopProcessID == "" || request.Token == "" {
		c.JSON(http.StatusBadRequest, BancardConfirmationResponse{
			Status:  StatusError,
			Message: "shop_process_id y token son obligatorios",
		})
		return
	}

	// Simular confirmación exitosa
	response := BancardConfirmationResponse{
		Status:              StatusSuccess,
		Message:             "Transacción confirmada exitosamente",
		TransactionID:       generateTransactionID(),
		AuthorizationNumber: generateAuthNumber(),
		TicketNumber:        generateTicketNumber(),
		ResponseCode:        "00",
		ResponseDescription: "Transacción aprobada",
		Amount:              "100000",
		Currency:            CurrencyPYG,
		Security: BancardSecurityInfo{
			Customer: BancardCustomerInfo{
				Document:     "12345678",
				DocumentType: "CI",
				Email:        "test@example.com",
				CellPhone:    "+595981234567",
			},
			CardInfo: BancardCardInfo{
				Bin:           "450000",
				Last4:         "1234",
				Brand:         BrandVisa,
				Type:          CardTypeCredit,
				Issuer:        "Banco Test",
				IssuerCountry: "PY",
			},
			RiskAnalysis: BancardRiskAnalysis{
				Score:          85,
				Recommendation: "APPROVE",
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// handleGetTransaction maneja la consulta de una transacción por ID
func (p *BancardPlugin) handleGetTransaction(c *gin.Context) {
	processID := c.Param("process_id")

	if processID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  StatusError,
			"message": "process_id es obligatorio",
		})
		return
	}

	// Simular información de transacción
	response := gin.H{
		"status":         StatusSuccess,
		"process_id":     processID,
		"transaction_id": generateTransactionID(),
		"amount":         "100000",
		"currency":       CurrencyPYG,
		"state":          "confirmed",
		"created_at":     "2025-09-16T12:30:45Z",
	}

	c.JSON(http.StatusOK, response)
}

// handleCheckout maneja la página de checkout
func (p *BancardPlugin) handleCheckout(c *gin.Context) {
	processID := c.Param("process_id")

	checkoutData := BancardCheckoutData{
		ProcessID:     processID,
		Amount:        "100000",
		Currency:      CurrencyPYG,
		ShopProcessID: processID,
		ReturnURL:     c.Query("return_url"),
		CancelURL:     c.Query("cancel_url"),
		OrderDetails: BancardOrderDetails{
			Amount:      "100000",
			Currency:    CurrencyPYG,
			Description: "Compra de ejemplo - Bancard VPOS",
		},
	}

	c.HTML(http.StatusOK, "bancard_checkout.html", gin.H{
		"data": checkoutData,
	})
}

// handleReturn maneja la página de retorno exitoso
func (p *BancardPlugin) handleReturn(c *gin.Context) {
	transactionID := c.Query("transaction_id")
	status := c.Query("status")

	c.HTML(http.StatusOK, "bancard_result.html", gin.H{
		"status":         status,
		"transaction_id": transactionID,
		"result":         StatusSuccess,
		"message":        "Transacción completada exitosamente",
	})
}

// handleCancel maneja la página de cancelación
func (p *BancardPlugin) handleCancel(c *gin.Context) {
	c.HTML(http.StatusOK, "bancard_result.html", gin.H{
		"status":  StatusCancelled,
		"result":  StatusCancelled,
		"message": "Transacción cancelada por el usuario",
	})
}

// handleEmulatorPayment maneja pagos del emulador
func (p *BancardPlugin) handleEmulatorPayment(c *gin.Context) {
	processID := c.Param("process_id")
	result := c.Query("result")

	var status string
	var message string
	var transactionID string

	switch result {
	case StatusSuccess:
		status = StatusSuccess
		message = "Pago procesado exitosamente"
		transactionID = generateTransactionID()
	case StatusError:
		status = StatusError
		message = "Error al procesar el pago"
	default:
		status = StatusCancelled
		message = "Pago cancelado"
	}

	simulationResult := BancardSimulationResult{
		Status:        status,
		ProcessID:     processID,
		Message:       message,
		TransactionID: transactionID,
	}

	c.JSON(http.StatusOK, simulationResult)
}

// handleEmulatorResult maneja el resultado del emulador
func (p *BancardPlugin) handleEmulatorResult(c *gin.Context) {
	result := c.Query("result")
	processID := c.Query("process_id")

	c.HTML(http.StatusOK, "bancard_result.html", gin.H{
		"result":     result,
		"process_id": processID,
		"status":     result,
	})
}

// Funciones auxiliares

func generateProcessID() string {
	return fmt.Sprintf("proc_%d_%d", rand.Int63(), rand.Intn(1000))
}

func generateTransactionID() string {
	return fmt.Sprintf("txn_%d", rand.Int63())
}

func generateAuthNumber() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func generateTicketNumber() string {
	return fmt.Sprintf("TKT%08d", rand.Intn(100000000))
}
