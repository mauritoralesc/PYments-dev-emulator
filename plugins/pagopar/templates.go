package pagopar

// GetPagoparTemplates devuelve todos los templates específicos de Pagopar
func GetPagoparTemplates() map[string]string {
	return map[string]string{
		"pagopar_checkout.html":  pagoparCheckoutHTML,
		"pagopar_result.html":    pagoparResultHTML,
		"pagopar_docs.html":      pagoparDocsHTML,
		"webhook_simulator.html": webhookSimulatorHTML,
	}
}

// Template para el checkout de Pagopar
const pagoparCheckoutHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Pagopar - Checkout</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .method { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 6px; border: 1px solid #ddd; cursor: pointer; }
        .method:hover { background: #e9ecef; }
        .selected { border-color: #007bff; background: #e7f3ff; }
        button { background: #007bff; color: white; border: none; padding: 12px 24px; border-radius: 4px; cursor: pointer; font-size: 16px; }
        button:hover { background: #0056b3; }
        .info { background: #d4edda; padding: 15px; border-radius: 4px; margin: 20px 0; border: 1px solid #c3e6cb; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏦 Pagopar - Checkout</h1>
        <div class="info">
            <p><strong>Hash del pedido:</strong> {{.hash}}</p>
            <p><strong>Forma de pago seleccionada:</strong> {{.formaPago}}</p>
        </div>
        
        <h2>Seleccionar método de pago:</h2>
        {{range .methods}}
        <div class="method" onclick="selectMethod('{{.FormaPago}}')">
            <h3>{{.Titulo}}</h3>
            <p>{{.Descripcion}}</p>
            <small>Comisión: {{.PorcentajeComision}}% - Mínimo: Gs. {{.MontoMinimo}}</small>
        </div>
        {{end}}
        
        <div style="margin-top: 30px;">
            <button onclick="processPayment('success')">✅ Simular Pago Exitoso</button>
            <button onclick="processPayment('error')" style="background: #dc3545; margin-left: 10px;">❌ Simular Error</button>
            <button onclick="processPayment('pending')" style="background: #ffc107; color: black; margin-left: 10px;">⏳ Simular Pendiente</button>
        </div>
        
        <div style="margin-top: 20px; border-top: 1px solid #ddd; padding-top: 20px;">
            <h3>🔄 Simular Flujo Completo</h3>
            <p><small>Esto simulará el flujo completo con redirect a tu aplicación</small></p>
            <button onclick="processPaymentWithRedirect('success')" style="background: #28a745;">✅ Pagar y Redirigir (Éxito)</button>
            <button onclick="processPaymentWithRedirect('error')" style="background: #dc3545; margin-left: 10px;">❌ Pagar y Redirigir (Error)</button>
        </div>
    </div>
    
    <script>
        function selectMethod(id) {
            document.querySelectorAll('.method').forEach(m => m.classList.remove('selected'));
            event.target.closest('.method').classList.add('selected');
        }
        
        function processPayment(result) {
            alert('Simulando resultado: ' + result);
            fetch('/emulator/webhook/{{.hash}}?result=' + result, {
                method: 'POST'
            }).then(response => response.json())
            .then(data => {
                console.log('Webhook simulado:', data);
                window.location.href = '/emulator/result?hash={{.hash}}&result=' + result;
            });
        }
        
        function processPaymentWithRedirect(result) {
            alert('Procesando pago y redirigiendo a la aplicación...');
            setTimeout(() => {
                window.location.href = '/resultado/{{.hash}}';
            }, 2000);
        }
    </script>
</body>
</html>`

// Template para el resultado de Pagopar
const pagoparResultHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Pagopar - Resultado</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; text-align: center; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .status { font-size: 48px; margin: 20px 0; }
        .message { font-size: 18px; margin: 20px 0; font-weight: bold; }
        .hash { background: #f8f9fa; padding: 15px; border-radius: 4px; font-family: monospace; word-break: break-all; margin: 20px 0; }
        button { background: #007bff; color: white; border: none; padding: 12px 24px; border-radius: 4px; cursor: pointer; font-size: 16px; margin: 10px; }
    </style>
</head>
<body>
    <div class="container">
        {{if eq .result "success"}}
        <div class="status">✅</div>
        <div class="message" style="color: #28a745;">Pago Exitoso</div>
        {{else if eq .result "error"}}
        <div class="status">❌</div>
        <div class="message" style="color: #dc3545;">Error en el Pago</div>
        {{else}}
        <div class="status">⏳</div>
        <div class="message" style="color: #ffc107;">Pago Pendiente</div>
        {{end}}
        
        <div class="hash">
            <strong>Hash:</strong><br>{{.hash}}
        </div>
        
        <button onclick="window.history.back()">← Volver</button>
        <button onclick="window.close()">Cerrar</button>
    </div>
</body>
</html>`

// Template para documentación de Pagopar
const pagoparDocsHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Pagopar - Documentación</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .route { background: #f8f9fa; padding: 15px; margin: 10px 0; border-radius: 6px; border-left: 4px solid #28a745; }
        .method { background: #007bff; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .path { font-family: monospace; font-size: 16px; margin: 10px 0; }
        h1 { color: #333; }
        .header { text-align: center; margin-bottom: 40px; }
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
        {{range .plugin.Routes}}
        <div class="route">
            <span class="method">{{.Method}}</span>
            <div class="path">{{.Path}}</div>
            <p><strong>Respuesta:</strong> {{.ResponseType}}</p>
        </div>
        {{end}}
        
        <h2>🧪 Ejemplo de Uso</h2>
        <div class="route">
            <p>Para probar este plugin, realiza una petición POST a:</p>
            <div class="path">http://localhost:{{.plugin.Port}}/api/comercios/2.0/iniciar-transaccion</div>
            <p>El emulador creará una orden y devolverá el hash para continuar con el flujo de checkout.</p>
        </div>
    </div>
</body>
</html>`

// Template para simulador de webhook
const webhookSimulatorHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Pagopar - Simulador de Webhook</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; background: #f8f9fa; padding: 40px 20px; }
        .container { max-width: 700px; margin: 0 auto; background: white; padding: 30px; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .webhook-data { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; font-family: monospace; font-size: 14px; white-space: pre-wrap; }
        .actions { text-align: center; margin-top: 30px; }
        .button { display: inline-block; padding: 15px 25px; margin: 10px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; font-weight: bold; }
        .success { background: #28a745; color: white; }
        .error { background: #dc3545; color: white; }
        .pending { background: #ffc107; color: #212529; }
        .cancel { background: #6c757d; color: white; }
        .info { background: #e7f3ff; padding: 15px; border-radius: 8px; margin-bottom: 20px; border-left: 4px solid #007bff; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔔 Simulador de Webhook Pagopar</h1>
            <p>Simula las notificaciones que Pagopar envía a tu sistema</p>
        </div>

        <div class="info">
            <strong>📝 Instrucciones:</strong>
            <p>Selecciona el resultado que deseas simular. Esto generará el webhook que Pagopar enviaría a tu endpoint de notificaciones.</p>
        </div>

        <div class="actions">
            <button class="button success" onclick="simulateWebhook('success')">
                ✅ Pago Exitoso
            </button>
            <button class="button error" onclick="simulateWebhook('error')">
                ❌ Pago Fallido  
            </button>
            <button class="button pending" onclick="simulateWebhook('pending')">
                ⏳ Pago Pendiente
            </button>
            <button class="button cancel" onclick="simulateWebhook('cancel')">
                🚫 Pago Cancelado
            </button>
        </div>
    </div>

    <script>
        function simulateWebhook(result) {
            alert('Simulando webhook: ' + result);
        }
    </script>
</body>
</html>`
