# Payment Emulator - Makefile

# Variables
BINARY_NAME=payment-emulator
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Colores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean install deps dev release

help: ## Mostrar ayuda
	@echo "${BLUE}Payment Emulator - Paraguay${NC}"
	@echo ""
	@echo "Comandos disponibles:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Compilar el binario
	@echo "${YELLOW}Compilando ${BINARY_NAME}...${NC}"
	go build ${LDFLAGS} -o ${BINARY_NAME}
	@echo "${GREEN}✅ Compilación exitosa${NC}"

run: build ## Compilar e iniciar el emulador
	@echo "${YELLOW}Iniciando Payment Emulator...${NC}"
	./${BINARY_NAME} start

dev: ## Modo desarrollo con hot reload
	@echo "${YELLOW}Iniciando en modo desarrollo...${NC}"
	go run main.go start --verbose

test: ## Ejecutar tests
	@echo "${YELLOW}Ejecutando tests...${NC}"
	go test -v ./...

clean: ## Limpiar archivos de compilación
	@echo "${YELLOW}Limpiando archivos...${NC}"
	rm -f ${BINARY_NAME}
	rm -rf ${BUILD_DIR}
	@echo "${GREEN}✅ Limpieza completada${NC}"

install: build ## Instalar globalmente
	@echo "${YELLOW}Instalando globalmente...${NC}"
	sudo cp ${BINARY_NAME} /usr/local/bin/
	@echo "${GREEN}✅ Instalado en /usr/local/bin/${BINARY_NAME}${NC}"

deps: ## Instalar dependencias
	@echo "${YELLOW}Instalando dependencias...${NC}"
	go mod download
	go mod tidy
	@echo "${GREEN}✅ Dependencias instaladas${NC}"

# Compilación cross-platform
release: clean ## Compilar para todas las plataformas
	@echo "${YELLOW}Compilando para múltiples plataformas...${NC}"
	@mkdir -p ${BUILD_DIR}
	
	@echo "  • Windows amd64..."
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe
	
	@echo "  • Linux amd64..."
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64
	
	@echo "  • macOS amd64..."
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64
	
	@echo "  • macOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64
	
	@echo "${GREEN}✅ Compilación cross-platform completada en ${BUILD_DIR}/${NC}"

# Comandos de utilidad
plugins: build ## Listar plugins disponibles
	./${BINARY_NAME} plugins list

add-plugin: ## Crear nuevo plugin (uso: make add-plugin NAME=miplugin)
	@if [ -z "$(NAME)" ]; then \
		echo "${RED}Error: Especifica NAME=nombre_plugin${NC}"; \
		exit 1; \
	fi
	./${BINARY_NAME} plugins add $(NAME)

demo: build ## Iniciar con configuración de demo
	@echo "${YELLOW}Iniciando demo...${NC}"
	@echo "Dashboard: ${BLUE}http://localhost:8000${NC}"
	@echo "Demo HTML: ${BLUE}file://$(PWD)/demo.html${NC}"
	@echo "Bancard:   ${BLUE}http://localhost:8001${NC}"
	@echo "Pagopar:   ${BLUE}http://localhost:8002${NC}"
	@echo ""
	./${BINARY_NAME} start

# Comandos de desarrollo
fmt: ## Formatear código
	go fmt ./...

lint: ## Ejecutar linter
	golangci-lint run

check: fmt lint test ## Ejecutar todas las verificaciones

# Docker (opcional)
docker-build: ## Compilar imagen Docker
	docker build -t payment-emulator .

docker-run: docker-build ## Ejecutar en Docker
	docker run -p 8000:8000 -p 8001:8001 -p 8002:8002 payment-emulator

# Información del proyecto
version: ## Mostrar versión
	@echo "${GREEN}Version: ${VERSION}${NC}"

info: ## Mostrar información del proyecto
	@echo "${BLUE}Payment Emulator - Paraguay${NC}"
	@echo "Version: ${VERSION}"
	@echo "Build: $(shell go version)"
	@echo "Platform: $(shell go env GOOS)/$(shell go env GOARCH)"
	@echo ""
	@echo "Archivos principales:"
	@find . -name "*.go" -not -path "./vendor/*" | head -10
	@echo ""
	@echo "Dependencias:"
	@go list -m all | head -5