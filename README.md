# PYment Dev Emulator - Paraguay

**Emulador de Pasarelas de Pago para Paraguay**

Un emulador completo de medios de pago paraguayos para desarrollo y testing. Soporta Bancard, Pagopar y permite agregar plugins personalizados.

## Características

- **CLI Completo**: Gestión mediante comandos con Cobra
- **Sistema de Plugins**: Cada medio de pago en su propio puerto
- **Emulación Realista**: Simula flujos reales de las pasarelas
- **Sin Docker**: Binario único sin dependencias
- **Cross-Platform**: Windows, Linux, macOS
- **Templates Embebidos**: No requiere archivos externos
- **Documentación Automática**: Cada plugin tiene su documentación

## Plugins Incluidos

### Bancard VPOS
- **Puerto**: 8001
- **Tipo**: iframe
- **Rutas**:
  - `POST /vpos/api/0.3/single_buy` - Iniciar pago
  - `POST /vpos/api/0.3/confirmation` - Confirmación
  - `POST /vpos/api/0.3/refund` - Reembolso

### Pagopar
- **Puerto**: 8002  
- **Tipo**: redirect
- **Rutas**:
  - `POST /api/comercios/2.0/iniciar-transaccion` - Crear orden
  - `POST /api/forma-pago/1.1/traer` - Listar métodos de pago
  - `POST /api/pedidos/1.1/traer` - Consultar estado
  - `GET /pagos/{hash}` - Página de checkout
  - `POST /emulator/webhook/{hash}` - Simulador de webhook

## Instalación

```bash
# Compilar el proyecto
go build -o payment-emulator

# O instalar globalmente
go install
```

## Uso

### Iniciar el Emulador

```bash
# Iniciar con configuración por defecto
./payment-emulator start

# Personalizar puerto y plugins
./payment-emulator start --port 9000 --plugins bancard,pagopar

# Solo dashboard sin plugins
./payment-emulator start --dashboard --plugins ""
```

### Gestionar Plugins

```bash
# Listar plugins disponibles
./payment-emulator plugins list

# Crear nuevo plugin
./payment-emulator plugins add miplugin
```

## Acceso

Una vez iniciado:

- **Dashboard Principal**: http://localhost:8000
- **Plugin Bancard**: http://localhost:8001  
- **Plugin Pagopar**: http://localhost:8002
- **API de Estado**: http://localhost:8000/api/plugins
- **Health Check**: http://localhost:8000/health

## Integración con tu App

### Ejemplo Bancard (iframe)

```html
<iframe src="http://localhost:8001/vpos/api/0.3/single_buy?amount=100000&return_url=https://tuapp.com/success&cancel_url=https://tuapp.com/cancel"></iframe>
```

### Ejemplo Pagopar (Flujo completo)

```javascript
// Step 1: Crear orden
const orderResponse = await fetch('http://localhost:8002/api/comercios/2.0/iniciar-transaccion', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    token: 'your_token',
    public_key: 'your_public_key',
    monto_total: 100000,
    comprador: { email: 'test@example.com', nombre: 'Juan' },
    compras_items: [{ nombre: 'Producto', precio_total: 100000 }]
  })
});

const order = await orderResponse.json();
const hash = order.resultado[0].data;

// Step 2: Redireccionar a checkout
window.location.href = `http://localhost:8002/pagos/${hash}`;

// Step 3: Webhook automático tras simulación
// Step 4: Consultar estado
const statusResponse = await fetch('http://localhost:8002/api/pedidos/1.1/traer', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    hash_pedido: hash,
    token: 'your_query_token',
    token_publico: 'your_public_key'
  })
});
```

## Flujo de Testing

1. **Integra tu app** con `localhost:8001` o `localhost:8002`
2. **Tu app abre** iframe/popup hacia el emulador
3. **Clickea resultado**: "Pago Exitoso", "Error", "Cancelado"  
4. **Emulador redirige** con parámetros apropiados
5. **Tu app procesa** la respuesta como en producción

## Crear Plugin Personalizado

```bash
# Crear plugin
./payment-emulator plugins add mipago

# Editar configuración
# plugins/mipago/config.yaml
```

Ejemplo de configuración:

```yaml
name: "Mi Pago"
description: "Mi sistema de pagos personalizado"
port: 8003
type: "iframe"  # o "redirección"
enabled: true
routes:
  - path: "/pay"
    method: "POST"
    response_type: "redirect"
  - path: "/webhook"
    method: "POST"
    response_type: "json"
```

## Compilación Cross-Platform

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o payment-emulator.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o payment-emulator-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o payment-emulator-macos
```

## Parámetros de Configuración

### Flags del comando `start`

- `--port, -p`: Puerto principal (default: 8000)
- `--plugins, -P`: Lista de plugins (default: bancard,pagopar)
- `--dashboard, -d`: Mostrar dashboard (default: true)

### Variables de Entorno

- `PAYMENT_EMULATOR_PORT`: Puerto por defecto
- `PAYMENT_EMULATOR_VERBOSE`: Modo verbose

### Testing

**Flujo Pagopar (4 pasos):**
1. **Crear orden**: `curl -X POST localhost:8002/api/comercios/2.0/iniciar-transaccion`
2. **Ir a checkout**: Abrir `localhost:8002/pagos/{hash}` en navegador
3. **Simular webhook**: Seleccionar resultado en interfaz
4. **Consultar estado**: `curl -X POST localhost:8002/api/pedidos/1.1/traer`

**Flujo Bancard (iframe):**
1. Integrar iframe: `localhost:8001/vpos/api/0.3/single_buy`
2. Simular resultado en interfaz
3. Procesar redirección

## Contribuir

1. Fork el proyecto
2. Crea tu feature branch
3. Agrega tus cambios  
4. Ejecuta las pruebas
5. Crea un Pull Request

## Licencia

MIT License - ver archivo LICENSE para detalles.

---

**Desarrollado para la comunidad paraguaya de developers** 🇵🇾