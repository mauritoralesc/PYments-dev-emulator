package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-emulator/internal/server"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia el emulador de medios de pago",
	Long: `Inicia el servidor principal y todos los plugins configurados.
Cada plugin se ejecutará en su puerto específico con su propia documentación.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer(cmd)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntP("port", "p", 8000, "Puerto principal del dashboard")
	startCmd.Flags().StringSliceP("plugins", "P", []string{"bancard", "pagopar"}, "Plugins a cargar")
	startCmd.Flags().BoolP("dashboard", "d", true, "Mostrar dashboard web")
}

func startServer(cmd *cobra.Command) {
	port, _ := cmd.Flags().GetInt("port")
	plugins, _ := cmd.Flags().GetStringSlice("plugins")
	dashboard, _ := cmd.Flags().GetBool("dashboard")

	fmt.Printf("Iniciando PYment Dev Emulator...\n")
	fmt.Printf("Dashboard: http://localhost:%d\n", port)
	fmt.Printf("API Docs: http://localhost:%d/docs\n", port)

	// Crear servidor principal
	mainServer := server.NewMainServer(port, dashboard)

	// Cargar plugins
	pluginServers := make([]*http.Server, 0)
	for i, pluginName := range plugins {
		pluginPort := port + i + 1
		pluginServer := server.NewPluginServer(pluginName, pluginPort)
		pluginServers = append(pluginServers, pluginServer)

		fmt.Printf(" Plugin %s: http://localhost:%d\n", pluginName, pluginPort)

		// Iniciar plugin en goroutine
		go func(srv *http.Server) {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Plugin server error: %v", err)
			}
		}(pluginServer)
	}

	// Iniciar servidor principal
	go func() {
		if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Main server error: %v", err)
		}
	}()

	fmt.Printf("\n Todos los servicios iniciados correctamente\n")
	fmt.Printf("Presiona Ctrl+C para detener\n\n")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Printf("\n Deteniendo servicios...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Detener servidor principal
	if err := mainServer.Shutdown(ctx); err != nil {
		log.Printf("Main server shutdown error: %v", err)
	}

	// Detener plugins
	for _, srv := range pluginServers {
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Plugin server shutdown error: %v", err)
		}
	}

	fmt.Printf(" Servicios detenidos correctamente\n")
}
