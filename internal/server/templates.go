package server

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

func loadTemplates(r *gin.Engine) {
	// Templates comunes del sistema (no específicos de plugins)
	dashboardHTML := `<!DOCTYPE html>
<html>
<head>
    <title>{{.title}}</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .plugin { background: #f8f9fa; padding: 20px; margin: 10px 0; border-radius: 6px; border-left: 4px solid #007bff; }
        .status { color: #28a745; font-weight: bold; }
        a { color: #007bff; text-decoration: none; }
        a:hover { text-decoration: underline; }
        h1 { color: #333; margin-bottom: 30px; }
        .header { text-align: center; margin-bottom: 40px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏦 {{.title}}</h1>
            <p>Emulador de medios de pago para desarrollo</p>
        </div>
        
        <h2>📊 Servicios Activos</h2>
        {{range .ports}}
        <div class="plugin">
            <h3>Plugin en Puerto {{.}}</h3>
            <p class="status">● Estado: Activo</p>
            <p><a href="http://localhost:{{.}}" target="_blank">📖 Ver Documentación</a></p>
        </div>
        {{end}}
        
        <h2>📚 Documentación</h2>
        <div class="plugin">
            <p><a href="/api/plugins">📋 Ver API de Plugins</a></p>
            <p><a href="/health">💓 Health Check</a></p>
        </div>
    </div>
</body>
</html>`

	pluginDocsHTML := `<!DOCTYPE html>
<html>
<head>
    <title>{{.plugin.Name}} - Documentación</title>
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
            <div class="path">http://localhost:{{.plugin.Port}}{{if .plugin.Routes}}{{(index .plugin.Routes 0).Path}}{{end}}</div>
            <p>El emulador mostrará una interfaz para simular diferentes resultados de pago.</p>
        </div>
    </div>
</body>
</html>`

	iframeEmulatorHTML := `<!DOCTYPE html>
<html>
<head>
    <title>{{.plugin.Name}} - Emulador</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f8f9fa; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .button { display: inline-block; padding: 12px 24px; margin: 10px; border: none; border-radius: 6px; font-size: 16px; cursor: pointer; text-decoration: none; text-align: center; }
        .success { background: #28a745; color: white; }
        .error { background: #dc3545; color: white; }
        .cancel { background: #6c757d; color: white; }
        .info { background: #e9ecef; padding: 15px; border-radius: 6px; margin-bottom: 20px; }
        h1 { color: #333; text-align: center; }
        .actions { text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h1>💳 {{.plugin.Name}}</h1>
        
        <div class="info">
            <h3>📄 Información de la Transacción</h3>
            <p><strong>Ruta:</strong> {{.route.Path}}</p>
            <p><strong>Método:</strong> {{.route.Method}}</p>
            {{range $key, $value := .params}}
            <p><strong>{{$key}}:</strong> {{index $value 0}}</p>
            {{end}}
        </div>
        
        <div class="actions">
            <h3>🎮 Simular Resultado</h3>
            <button class="button success" onclick="simulateResult('success')">✅ Pago Exitoso</button>
            <button class="button error" onclick="simulateResult('error')">❌ Error de Pago</button>
            <button class="button cancel" onclick="simulateResult('cancel')">🚫 Cancelado por Usuario</button>
        </div>
    </div>

    <script>
        function simulateResult(result) {
            const params = new URLSearchParams(window.location.search);
            const returnUrl = params.get('return_url') || params.get('success_url') || 'about:blank';
            const cancelUrl = params.get('cancel_url') || params.get('error_url') || 'about:blank';
            
            let redirectUrl;
            switch(result) {
                case 'success':
                    redirectUrl = returnUrl + '?status=success&transaction_id=' + Math.random().toString(36).substr(2, 9);
                    break;
                case 'error':
                    redirectUrl = cancelUrl + '?status=error&error_code=E001&error_message=Payment%20declined';
                    break;
                case 'cancel':
                    redirectUrl = cancelUrl + '?status=cancelled';
                    break;
            }
            
            // Si estamos en iframe, usar parent
            if (window.parent !== window) {
                window.parent.location.href = redirectUrl;
            } else {
                window.location.href = redirectUrl;
            }
        }
    </script>
</body>
</html>`

	popupEmulatorHTML := `<!DOCTYPE html>
<html>
<head>
    <title>{{.plugin.Name}} - Emulador Popup</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f8f9fa; }
        .container { max-width: 450px; margin: 0 auto; background: white; padding: 25px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .button { display: inline-block; padding: 10px 20px; margin: 8px; border: none; border-radius: 6px; font-size: 14px; cursor: pointer; text-decoration: none; text-align: center; }
        .success { background: #28a745; color: white; }
        .error { background: #dc3545; color: white; }
        .cancel { background: #6c757d; color: white; }
        .info { background: #e9ecef; padding: 12px; border-radius: 6px; margin-bottom: 15px; font-size: 14px; }
        h1 { color: #333; text-align: center; font-size: 24px; }
        .actions { text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h1>💳 {{.plugin.Name}}</h1>
        
        <div class="info">
            <h3>📄 Información del Pago</h3>
            <p><strong>Ruta:</strong> {{.route.Path}}</p>
            {{range $key, $value := .params}}
            <p><strong>{{$key}}:</strong> {{index $value 0}}</p>
            {{end}}
        </div>
        
        <div class="actions">
            <h3>🎮 Simular Resultado</h3>
            <button class="button success" onclick="simulateResult('success')">✅ Confirmar Pago</button>
            <button class="button error" onclick="simulateResult('error')">❌ Error en Pago</button>
            <button class="button cancel" onclick="simulateResult('cancel')">🚫 Cancelar</button>
        </div>
    </div>

    <script>
        function simulateResult(result) {
            const params = new URLSearchParams(window.location.search);
            const callbackUrl = params.get('callback_url') || params.get('return_url') || window.opener;
            
            let resultData;
            switch(result) {
                case 'success':
                    resultData = {
                        status: 'success',
                        transaction_id: Math.random().toString(36).substr(2, 9),
                        amount: params.get('amount'),
                        currency: params.get('currency') || 'PYG'
                    };
                    break;
                case 'error':
                    resultData = {
                        status: 'error',
                        error_code: 'E001',
                        error_message: 'Payment declined'
                    };
                    break;
                case 'cancel':
                    resultData = {
                        status: 'cancelled'
                    };
                    break;
            }
            
            // Si es popup, enviar mensaje al padre
            if (window.opener) {
                window.opener.postMessage(resultData, '*');
                window.close();
            } else if (callbackUrl) {
                const url = new URL(callbackUrl);
                Object.keys(resultData).forEach(key => {
                    url.searchParams.set(key, resultData[key]);
                });
                window.location.href = url.toString();
            }
        }
    </script>
</body>
</html>`

	// Registrar templates comunes
	templ := template.New("")
	templ = template.Must(templ.New("dashboard.html").Parse(dashboardHTML))
	templ = template.Must(templ.New("plugin_docs.html").Parse(pluginDocsHTML))
	templ = template.Must(templ.New("iframe_emulator.html").Parse(iframeEmulatorHTML))
	templ = template.Must(templ.New("popup_emulator.html").Parse(popupEmulatorHTML))

	// Agregar templates específicos de plugins
	loadPluginSpecificTemplates(templ)

	r.SetHTMLTemplate(templ)
}

// loadPluginSpecificTemplates carga templates específicos de cada plugin
func loadPluginSpecificTemplates(templ *template.Template) {
	// Templates de Pagopar
	pagoparTemplates := map[string]string{
		"pagopar_checkout.html": `<!DOCTYPE html>
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
</html>`,

		"pagopar_result.html": `<!DOCTYPE html>
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
</html>`,
	}

	// Templates de Bancard
	bancardTemplates := map[string]string{
		"bancard_checkout.html": `<!DOCTYPE html>
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
        .actions { text-align: center; margin-top: 30px; }
        .button { display: inline-block; padding: 15px 30px; margin: 10px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; font-weight: bold; }
        .primary { background: #1e3c72; color: white; }
        .error { background: #dc3545; color: white; }
        .secondary { background: #6c757d; color: white; }
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
        function processPayment(result) {
            const processId = '{{.data.ProcessID}}';
            
            if (result === 'success') {
                document.querySelector('.primary').innerHTML = '⏳ Procesando...';
                document.querySelector('.primary').disabled = true;
                
                setTimeout(() => {
                    const returnUrl = '{{.data.ReturnURL}}' || '/bancard/return';
                    const transactionId = 'TXN' + Math.random().toString(36).substr(2, 9).toUpperCase();
                    window.location.href = returnUrl + '?status=success&transaction_id=' + transactionId + '&process_id=' + processId;
                }, 2000);
            } else if (result === 'error') {
                alert('Simulando error: Tarjeta rechazada');
                const returnUrl = '{{.data.ReturnURL}}' || '/bancard/return';
                window.location.href = returnUrl + '?status=error&error_code=05&process_id=' + processId;
            } else {
                const cancelUrl = '{{.data.CancelURL}}' || '/bancard/cancel';
                window.location.href = cancelUrl + '?status=cancelled&process_id=' + processId;
            }
        }
    </script>
</body>
</html>`,

		"bancard_result.html": `<!DOCTYPE html>
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
        .button { display: inline-block; padding: 15px 30px; margin: 10px; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; text-decoration: none; font-weight: bold; }
        .primary { background: #1e3c72; color: white; }
        .secondary { background: #6c757d; color: white; }
    </style>
</head>
<body>
    <div class="container">
        {{if eq .result "success"}}
        <div class="status">✅</div>
        <div class="message" style="color: #28a745;">¡Pago Exitoso!</div>
        {{else if eq .result "error"}}
        <div class="status">❌</div>
        <div class="message" style="color: #dc3545;">Error en el Pago</div>
        {{else}}
        <div class="status">🚫</div>
        <div class="message" style="color: #ffc107;">Pago Cancelado</div>
        {{end}}
        
        {{if .transaction_id}}
        <div class="details">
            <h4>📄 Detalles de la Transacción</h4>
            <p><strong>ID de Transacción:</strong> {{.transaction_id}}</p>
            <p><strong>Estado:</strong> {{.status}}</p>
            {{if .process_id}}<p><strong>Process ID:</strong> {{.process_id}}</p>{{end}}
        </div>
        {{end}}
        
        <div style="margin-top: 30px;">
            <button class="button primary" onclick="window.history.back()">← Volver</button>
            <button class="button secondary" onclick="window.close()">Cerrar</button>
        </div>
    </div>
</body>
</html>`,
	}

	// Registrar templates de Pagopar
	for name, content := range pagoparTemplates {
		templ = template.Must(templ.New(name).Parse(content))
	}

	// Registrar templates de Bancard
	for name, content := range bancardTemplates {
		templ = template.Must(templ.New(name).Parse(content))
	}
}
