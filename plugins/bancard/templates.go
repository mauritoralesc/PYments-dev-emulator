package bancard

// GetBancardTemplates devuelve todos los templates específicos de Bancard
func GetBancardTemplates() map[string]string {
	return map[string]string{
		"bancard_checkout.html": bancardCheckoutHTML,
		"bancard_result.html":   bancardResultHTML,
		"bancard_docs.html":     bancardDocsHTML,
	}
}

// Template para el checkout de Bancard
const bancardCheckoutHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Bancard VPOS - Checkout</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%); min-height: 100vh; }
        .container { max-width: 500px; margin: 0 auto; padding: 40px 20px; }
        .checkout-card { background: white; border-radius: 12px; padding: 30px; box-shadow: 0 10px 30px rgba(0,0,0,0.2); }
        .logo { text-align: center; margin-bottom: 30px; }
        .logo h1 { color: #1e3c72; margin: 0; font-size: 28px; }
        .order-info { background: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 30px; }
        .card-form { margin: 20px 0; }
        .form-group { margin-bottom: 15px; }
        .form-group label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
        .form-group input { width: 100%; padding: 12px; border: 1px solid #ddd; border-radius: 6px; font-size: 16px; }
        .form-row { display: flex; gap: 15px; }
        .form-row .form-group { flex: 1; }
        .actions { text-align: center; margin-top: 30px; }
        .button { display: inline-block; padding: 15px 30px; margin: 10px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; font-weight: bold; }
        .primary { background: #1e3c72; color: white; }
        .secondary { background: #6c757d; color: white; }
        .error { background: #dc3545; color: white; }
        .primary:hover { background: #2a5298; }
        .security-info { background: #e7f3ff; padding: 15px; border-radius: 8px; margin-top: 20px; border-left: 4px solid #007bff; }
    </style>
</head>
<body>
    <div class="container">
        <div class="checkout-card">
            <div class="logo">
                <h1>💳 Bancard VPOS</h1>
                <p>Procesamiento Seguro de Pagos</p>
            </div>
            
            <div class="order-info">
                <h3>📋 Información de la Transacción</h3>
                <p><strong>Process ID:</strong> {{.data.ProcessID}}</p>
                <p><strong>Monto:</strong> ₲ {{.data.Amount}}</p>
                <p><strong>Moneda:</strong> {{.data.Currency}}</p>
                <p><strong>Descripción:</strong> {{.data.OrderDetails.Description}}</p>
            </div>
            
            <div class="card-form">
                <h3>💳 Información de la Tarjeta</h3>
                <div class="form-group">
                    <label for="card-number">Número de Tarjeta</label>
                    <input type="text" id="card-number" placeholder="1234 5678 9012 3456" maxlength="19" 
                           oninput="formatCardNumber(this)">
                </div>
                
                <div class="form-row">
                    <div class="form-group">
                        <label for="expiry">Vencimiento</label>
                        <input type="text" id="expiry" placeholder="MM/AA" maxlength="5" 
                               oninput="formatExpiry(this)">
                    </div>
                    <div class="form-group">
                        <label for="cvv">CVV</label>
                        <input type="text" id="cvv" placeholder="123" maxlength="4">
                    </div>
                </div>
                
                <div class="form-group">
                    <label for="cardholder">Titular de la Tarjeta</label>
                    <input type="text" id="cardholder" placeholder="JUAN PEREZ" style="text-transform: uppercase;">
                </div>
                
                <div class="form-group">
                    <label for="document">Documento</label>
                    <input type="text" id="document" placeholder="12345678">
                </div>
            </div>
            
            <div class="security-info">
                <small>🔒 <strong>Transacción Segura:</strong> Sus datos están protegidos con encriptación SSL de 256 bits.</small>
            </div>
            
            <div class="actions">
                <button class="button primary" onclick="processPayment('success')">
                    ✅ Procesar Pago (₲ {{.data.Amount}})
                </button>
                <button class="button error" onclick="processPayment('error')">
                    ❌ Simular Error
                </button>
                <button class="button secondary" onclick="processPayment('cancel')">
                    🚫 Cancelar
                </button>
            </div>
        </div>
    </div>

    <script>
        function formatCardNumber(input) {
            let value = input.value.replace(/\s/g, '').replace(/[^0-9]/gi, '');
            let formattedValue = value.match(/.{1,4}/g)?.join(' ') || value;
            input.value = formattedValue;
        }

        function formatExpiry(input) {
            let value = input.value.replace(/\D/g, '');
            if (value.length >= 2) {
                value = value.substring(0, 2) + '/' + value.substring(2, 4);
            }
            input.value = value;
        }

        function processPayment(result) {
            const processId = '{{.data.ProcessID}}';
            
            if (result === 'success') {
                // Simular procesamiento
                document.querySelector('.primary').innerHTML = '⏳ Procesando...';
                document.querySelector('.primary').disabled = true;
                
                setTimeout(() => {
                    // Simular respuesta exitosa
                    const returnUrl = '{{.data.ReturnURL}}' || '/bancard/return';
                    const transactionId = 'TXN' + Math.random().toString(36).substr(2, 9).toUpperCase();
                    window.location.href = returnUrl + '?status=success&transaction_id=' + transactionId + '&process_id=' + processId;
                }, 2000);
            } else if (result === 'error') {
                alert('Simulando error: Tarjeta rechazada por el banco emisor');
                const returnUrl = '{{.data.ReturnURL}}' || '/bancard/return';
                window.location.href = returnUrl + '?status=error&error_code=05&process_id=' + processId;
            } else {
                const cancelUrl = '{{.data.CancelURL}}' || '/bancard/cancel';
                window.location.href = cancelUrl + '?status=cancelled&process_id=' + processId;
            }
        }
    </script>
</body>
</html>`

// Template para el resultado de Bancard
const bancardResultHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Bancard VPOS - Resultado</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; background: #f8f9fa; padding: 40px 20px; text-align: center; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.1); }
        .status { font-size: 64px; margin: 20px 0; }
        .message { font-size: 20px; margin: 20px 0; font-weight: bold; }
        .details { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; text-align: left; }
        .details p { margin: 5px 0; }
        .button { display: inline-block; padding: 15px 30px; margin: 10px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; font-weight: bold; }
        .primary { background: #1e3c72; color: white; }
        .secondary { background: #6c757d; color: white; }
    </style>
</head>
<body>
    <div class="container">
        {{if eq .result "success"}}
        <div class="status">Éxito</div>
        <div class="message" style="color: #28a745;">¡Pago Exitoso!</div>
        <p>Su transacción ha sido procesada correctamente.</p>
        {{else if eq .result "error"}}
        <div class="status">Error</div>
        <div class="message" style="color: #dc3545;">Error en el Pago</div>
        <p>No se pudo procesar su pago. Por favor, intente nuevamente.</p>
        {{else}}
        <div class="status">Pago Cancelado</div>
        <div class="message" style="color: #ffc107;">Pago Cancelado</div>
        <p>La transacción fue cancelada por el usuario.</p>
        {{end}}
        
        {{if .transaction_id}}
        <div class="details">
            <h4>📄 Detalles de la Transacción</h4>
            <p><strong>ID de Transacción:</strong> {{.transaction_id}}</p>
            <p><strong>Estado:</strong> {{.status}}</p>
            {{if .process_id}}<p><strong>Process ID:</strong> {{.process_id}}</p>{{end}}
            <p><strong>Fecha:</strong> <span id="current-date"></span></p>
        </div>
        {{end}}
        
        <div style="margin-top: 30px;">
            <button class="button primary" onclick="window.history.back()">
                ← Volver
            </button>
            <button class="button secondary" onclick="window.close()">
                Cerrar
            </button>
        </div>
    </div>

    <script>
        document.getElementById('current-date').textContent = new Date().toLocaleString('es-PY');
    </script>
</body>
</html>`

// Template para documentación de Bancard
const bancardDocsHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Bancard VPOS - Documentación</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .route { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 6px; border-left: 4px solid #1e3c72; }
        .method { background: #1e3c72; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .path { font-family: monospace; font-size: 16px; margin: 10px 0; }
        h1 { color: #1e3c72; }
        .header { text-align: center; margin-bottom: 40px; }
        code { background: #f8f9fa; padding: 2px 6px; border-radius: 4px; font-family: monospace; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💳 {{.plugin.Name}}</h1>
            <p>{{.plugin.Description}}</p>
            <p><strong>Tipo:</strong> {{.plugin.Type}} | <strong>Puerto:</strong> {{.plugin.Port}}</p>
        </div>
        
        <h2>🛤️ Rutas Disponibles</h2>
        
        <div class="route">
            <span class="method">POST</span>
            <div class="path">/vpos/api/0.3/single_buy</div>
            <p><strong>Descripción:</strong> Crear una nueva transacción de compra simple</p>
            <p><strong>Respuesta:</strong> JSON con process_id y redirect_url</p>
        </div>
        
        <div class="route">
            <span class="method">POST</span>
            <div class="path">/vpos/api/0.3/confirmation</div>
            <p><strong>Descripción:</strong> Confirmar una transacción existente</p>
            <p><strong>Respuesta:</strong> JSON con detalles de confirmación</p>
        </div>
        
        <div class="route">
            <span class="method">GET</span>
            <div class="path">/bancard/checkout/:process_id</div>
            <p><strong>Descripción:</strong> Página de checkout para completar el pago</p>
            <p><strong>Respuesta:</strong> HTML con formulario de pago</p>
        </div>
        
        <h2>🧪 Ejemplo de Uso</h2>
        <div class="route">
            <p><strong>1. Crear transacción:</strong></p>
            <div class="path">POST http://localhost:{{.plugin.Port}}/vpos/api/0.3/single_buy</div>
            <p>Payload: <code>{"public_key": "key", "operation": {"token": "token", "amount": "100000", "shop_process_id": "123"}}</code></p>
            
            <p><strong>2. Redirigir al checkout:</strong></p>
            <div class="path">GET http://localhost:{{.plugin.Port}}/bancard/checkout/{process_id}</div>
            
            <p><strong>3. Confirmar transacción:</strong></p>
            <div class="path">POST http://localhost:{{.plugin.Port}}/vpos/api/0.3/confirmation</div>
        </div>
    </div>
</body>
</html>`
