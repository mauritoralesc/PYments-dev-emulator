package pagopar

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"payment-emulator/internal/plugins"

	"github.com/gin-gonic/gin"
)

// PagoparPlugin implementa PaymentPlugin para Pagopar
type PagoparPlugin struct {
	name       string
	pluginType string
	port       int
	config     *plugins.Plugin
}

// NewPagoparPlugin crea una nueva instancia del plugin de Pagopar
func NewPagoparPlugin(config *plugins.Plugin) plugins.PaymentPlugin {
	return &PagoparPlugin{
		name:       "Pagopar",
		pluginType: "popup",
		config:     config,
	}
}

// GetName devuelve el nombre del plugin
func (p *PagoparPlugin) GetName() string {
	return p.name
}

// GetType devuelve el tipo del plugin
func (p *PagoparPlugin) GetType() string {
	return p.pluginType
}

// SetupRoutes configura todas las rutas específicas de Pagopar
func (p *PagoparPlugin) SetupRoutes(r *gin.Engine) {
	p.setupAPIRoutes(r)
	p.setupCheckoutRoutes(r)
	p.setupEmulatorRoutes(r)
}

// GetTemplates devuelve los templates específicos de Pagopar
func (p *PagoparPlugin) GetTemplates() map[string]string {
	return GetPagoparTemplates()
}

// HandlePaymentRequest maneja peticiones de pago genéricas
func (p *PagoparPlugin) HandlePaymentRequest(c *gin.Context, route *plugins.Route) {
	// Para Pagopar, las peticiones son manejadas por rutas específicas
	// Aquí podríamos implementar un fallback o documentación
	c.HTML(http.StatusOK, "pagopar_docs.html", gin.H{
		"plugin": p.config,
		"route":  route,
		"params": c.Request.URL.Query(),
	})
}

// setupAPIRoutes configura las rutas de la API de Pagopar
func (p *PagoparPlugin) setupAPIRoutes(r *gin.Engine) {
	// Step 1: Iniciar transacción - Crear orden en Pagopar
	r.POST("/api/comercios/2.0/iniciar-transaccion", p.handleIniciarTransaccion)

	// Step 2: Obtener métodos de pago disponibles
	r.POST("/api/forma-pago/1.1/traer", p.handleGetPaymentMethods)
	r.POST("/api/forma-pago/1.1/traer/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/forma-pago/1.1/traer")
	})

	// Step 4: Consultar estado de pedido
	r.POST("/api/pedidos/1.1/traer", p.handleGetOrderStatus)
	r.POST("/getOrderStatus", p.handleGetOrderStatusLegacy)

	// Webhooks
	r.POST("/api/webhook/confirm", p.handleWebhookConfirm)
	r.POST("/api/webhook/reversal", p.handleWebhookReversal)
}

// setupCheckoutRoutes configura las rutas del checkout
func (p *PagoparPlugin) setupCheckoutRoutes(r *gin.Engine) {
	// Step 2: Página de checkout de Pagopar
	r.GET("/pagos/:hash", p.handleCheckout)

	// Endpoint para resultado
	r.GET("/resultado/:hash", p.handleResult)
}

// setupEmulatorRoutes configura las rutas del emulador
func (p *PagoparPlugin) setupEmulatorRoutes(r *gin.Engine) {
	// Simular webhook de notificación
	r.POST("/emulator/webhook/:hash", p.handleEmulatorWebhook)

	// Página de resultado del pago
	r.GET("/emulator/result", p.handleEmulatorResult)
}

// handleIniciarTransaccion maneja la creación de una nueva transacción
func (p *PagoparPlugin) handleIniciarTransaccion(c *gin.Context) {
	var request PagoparOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, PagoparOrderResponse{
			Respuesta: false,
			Error:     "Datos JSON inválidos: " + err.Error(),
		})
		return
	}

	// Validar campos obligatorios
	if request.Token == "" {
		c.JSON(http.StatusBadRequest, PagoparOrderResponse{
			Respuesta: false,
			Error:     "Campo obligatorio faltante: token",
		})
		return
	}

	if request.MontoTotal == "" {
		c.JSON(http.StatusBadRequest, PagoparOrderResponse{
			Respuesta: false,
			Error:     "Campo obligatorio faltante: monto_total",
		})
		return
	}

	// Generar respuesta de orden creada
	hash := generateOrderHash()
	orderNumber := generateOrderNumber()

	response := PagoparOrderResponse{
		Respuesta: true,
		Resultado: []PagoparOrderResult{
			{
				Data:   hash,
				Pedido: orderNumber,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}

// handleGetPaymentMethods maneja la obtención de métodos de pago
func (p *PagoparPlugin) handleGetPaymentMethods(c *gin.Context) {
	var request PagoparPaymentMethodsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, PagoparPaymentMethodsResponse{
			Respuesta: false,
		})
		return
	}

	// Validar tokens
	if request.Token == "" || request.TokenPublico == "" {
		c.JSON(http.StatusBadRequest, PagoparPaymentMethodsResponse{
			Respuesta: false,
		})
		return
	}

	response := PagoparPaymentMethodsResponse{
		Respuesta: true,
		Resultado: getPaymentMethods(),
	}

	c.JSON(http.StatusOK, response)
}

// handleGetOrderStatus maneja la consulta de estado de pedido
func (p *PagoparPlugin) handleGetOrderStatus(c *gin.Context) {
	var request PagoparOrderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, PagoparOrderStatusResponse{
			Respuesta: false,
		})
		return
	}

	// Validar campos obligatorios
	if request.HashPedido == "" || request.Token == "" || request.TokenPublico == "" {
		c.JSON(http.StatusBadRequest, PagoparOrderStatusResponse{
			Respuesta: false,
		})
		return
	}

	status := getOrderStatus(request.HashPedido)
	response := PagoparOrderStatusResponse{
		Respuesta: true,
		Resultado: []PagoparOrderStatusData{status},
	}

	c.JSON(http.StatusOK, response)
}

// handleGetOrderStatusLegacy maneja la consulta de estado (versión legacy)
func (p *PagoparPlugin) handleGetOrderStatusLegacy(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	hashPedido, hashOk := request["hash_pedido"].(string)
	token, tokenOk := request["token"].(string)

	if !hashOk || !tokenOk {
		c.JSON(http.StatusBadRequest, gin.H{
			"respuesta": false,
			"mensaje":   "Faltan campos requeridos: hash_pedido, token",
		})
		return
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"respuesta": false,
			"mensaje":   "Token inválido",
		})
		return
	}

	status := getOrderStatus(hashPedido)
	response := gin.H{
		"respuesta": true,
		"resultado": []PagoparOrderStatusData{status},
	}

	c.JSON(http.StatusOK, response)
}

// handleCheckout maneja la página de checkout
func (p *PagoparPlugin) handleCheckout(c *gin.Context) {
	hash := c.Param("hash")
	formaPago := c.Query("forma_pago")

	// Usar template específico de Pagopar
	c.HTML(http.StatusOK, "pagopar_checkout.html", gin.H{
		"hash":      hash,
		"formaPago": formaPago,
		"methods":   getPaymentMethods(),
	})
}

// handleResult maneja el resultado final
func (p *PagoparPlugin) handleResult(c *gin.Context) {
	hash := c.Param("hash")
	// Redirigir a la aplicación principal
	redirectURL := fmt.Sprintf("http://localhost/pr/%s", hash)
	c.Redirect(http.StatusFound, redirectURL)
}

// handleWebhookConfirm maneja confirmación de webhook
func (p *PagoparPlugin) handleWebhookConfirm(c *gin.Context) {
	var request map[string]interface{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	hashPedido := request["hash_pedido"].(string)
	confirmationData := generateWebhookData(hashPedido, PaymentStatusSuccess)

	response := PagoparWebhookData{
		Resultado: []PagoparOrderStatusData{confirmationData},
		Respuesta: true,
	}

	c.JSON(http.StatusOK, response)
}

// handleWebhookReversal maneja webhook de reversal
func (p *PagoparPlugin) handleWebhookReversal(c *gin.Context) {
	var request map[string]interface{}
	c.ShouldBindJSON(&request)

	hashPedido := request["hash_pedido"].(string)
	reversalData := generateWebhookData(hashPedido, PaymentStatusError)

	response := PagoparWebhookData{
		Resultado: []PagoparOrderStatusData{reversalData},
		Respuesta: true,
	}

	c.JSON(http.StatusOK, response)
}

// handleEmulatorWebhook maneja webhooks del emulador
func (p *PagoparPlugin) handleEmulatorWebhook(c *gin.Context) {
	hash := c.Param("hash")
	result := c.Query("result")

	webhookData := generateWebhookData(hash, result)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Simulador de webhook",
		"webhook_data": webhookData,
		"hash":         hash,
	})
}

// handleEmulatorResult maneja la página de resultado del emulador
func (p *PagoparPlugin) handleEmulatorResult(c *gin.Context) {
	hash := c.Query("hash")
	result := c.Query("result")

	if result == PaymentStatusSuccess {
		redirectURL := fmt.Sprintf("http://localhost/pr/%s", hash)
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	c.HTML(http.StatusOK, "pagopar_result.html", gin.H{
		"hash":   hash,
		"result": result,
	})
}

// Funciones auxiliares

func generateOrderHash() string {
	return "ad57c9c94f745fdd9bc9093bb409297607264af1a904e6300e71c24f15d6" + fmt.Sprintf("%03d", 100+rand.Intn(900))
}

func generateOrderNumber() string {
	return fmt.Sprintf("%d", 1750+rand.Intn(1000))
}

func generateToken(hash string) string {
	return fmt.Sprintf("token_%s", hash[:10])
}

func generateWebhookToken(hashPedido string) string {
	tokenSecret := "test_webhook_secret_123"
	data := tokenSecret + hashPedido
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func getPaymentMethods() []PagoparPaymentMethod {
	return []PagoparPaymentMethod{
		{
			FormaPago:          "25",
			Titulo:             "PIX",
			Descripcion:        "PIX vía QR",
			MontoMinimo:        "1000",
			PorcentajeComision: "3.00",
		},
		{
			FormaPago:          "24",
			Titulo:             "Pago QR",
			Descripcion:        "Pagá con la app de tu banco, financiera o cooperativa a través de un QR",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "18",
			Titulo:             "Zimple",
			Descripcion:        "Utilice sus fondos de Zimple",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:            "9",
			Titulo:               "Tarjetas de crédito",
			Descripcion:          "Acepta Visa, Mastercard, American Express, Cabal, Panal, Discover, Diners Club.",
			MontoMinimo:          "1000",
			PorcentajeComision:   "6.82",
			PagosInternacionales: false,
		},
		{
			FormaPago:          "10",
			Titulo:             "Tigo Money",
			Descripcion:        "Utilice sus fondos de Tigo Money",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "11",
			Titulo:             "Transferencia Bancaria",
			Descripcion:        "Pago con transferencias bancarias. Los pagos se procesan de 08:30 a 17:30 hs.",
			MontoMinimo:        "1000",
			PorcentajeComision: "3.30",
		},
		{
			FormaPago:          "12",
			Titulo:             "Billetera Personal",
			Descripcion:        "Utilice sus fondos de Billetera Personal",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "13",
			Titulo:             "Pago Móvil",
			Descripcion:        "Usando la App Pago Móvil / www.infonet.com.py",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "20",
			Titulo:             "Wally",
			Descripcion:        "Utilice sus fondos de Wally",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "23",
			Titulo:             "Giros Claro",
			Descripcion:        "Utilice sus fondos de Billetera Claro",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "22",
			Titulo:             "Wepa",
			Descripcion:        "Acercándose a las bocas de pagos habilitadas luego de confirmar el pedido",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "2",
			Titulo:             "Aqui Pago",
			Descripcion:        "Acercándose a las bocas de pagos habilitadas luego de confirmar el pedido",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "3",
			Titulo:             "Pago Express",
			Descripcion:        "Acercándose a las bocas de pagos habilitadas luego de confirmar el pedido",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
		{
			FormaPago:          "15",
			Titulo:             "Infonet Cobranzas",
			Descripcion:        "Acercándose a las bocas de pagos habilitadas luego de confirmar el pedido",
			MontoMinimo:        "1000",
			PorcentajeComision: "6.82",
		},
	}
}

func getOrderStatus(hash string) PagoparOrderStatusData {
	return PagoparOrderStatusData{
		Pagado:                   true,
		NumeroComprobanteInterno: "8230473",
		UltimoMensajeError:       nil,
		FormaPago:                "Tarjetas de crédito/débito",
		FechaPago:                "2025-09-16 12:30:45.123456",
		Monto:                    "100000.00",
		FechaMaximaPago:          "2025-09-23 14:14:48",
		HashPedido:               hash,
		NumeroPedido:             "1746",
		Cancelado:                false,
		FormaPagoIdentificador:   "1",
		Token:                    generateToken(hash),
		MensajeResultadoPago: map[string]interface{}{
			"titulo":      "Pago procesado exitosamente",
			"descripcion": "Comprobante: 8230473. Tu pago ha sido procesado correctamente.",
		},
	}
}

func generateWebhookData(hash, result string) PagoparOrderStatusData {
	isPaid := result == PaymentStatusSuccess
	var fechaPago interface{} = nil
	if isPaid {
		fechaPago = "2025-09-16 12:30:45.123456"
	}

	return PagoparOrderStatusData{
		Pagado:                   isPaid,
		NumeroComprobanteInterno: "8230473",
		UltimoMensajeError:       nil,
		FormaPago:                "Tarjetas de crédito/débito",
		FechaPago:                fechaPago,
		Monto:                    "100000.00",
		FechaMaximaPago:          "2025-09-23 14:14:48",
		HashPedido:               hash,
		NumeroPedido:             "1750",
		Cancelado:                result == PaymentStatusCancel,
		FormaPagoIdentificador:   "9",
		Token:                    generateToken(hash),
	}
}
