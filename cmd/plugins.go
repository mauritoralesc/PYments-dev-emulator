package cmd

import (
	"fmt"
	"payment-emulator/internal/plugins"

	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Gestionar plugins de medios de pago",
}

var listPluginsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los plugins disponibles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📦 Plugins disponibles:")

		availablePlugins := plugins.GetAvailablePlugins()
		for _, plugin := range availablePlugins {
			status := "❌ Deshabilitado"
			if plugin.Enabled {
				status = "✅ Habilitado"
			}
			fmt.Printf("  • %s - %s %s\n", plugin.Name, plugin.Description, status)
			fmt.Printf("    Puerto: %d | Tipo: %s\n", plugin.Port, plugin.Type)
			fmt.Printf("    Rutas: %v\n\n", plugin.Routes)
		}
	},
}

var addPluginCmd = &cobra.Command{
	Use:   "add [nombre]",
	Short: "Agrega un nuevo plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		fmt.Printf("➕ Creando plugin: %s\n", pluginName)

		err := plugins.CreatePluginTemplate(pluginName)
		if err != nil {
			fmt.Printf("❌ Error creando plugin: %v\n", err)
			return
		}

		fmt.Printf("✅ Plugin %s creado exitosamente\n", pluginName)
		fmt.Printf("📁 Edita el archivo: plugins/%s/config.yaml\n", pluginName)
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
	pluginsCmd.AddCommand(listPluginsCmd)
	pluginsCmd.AddCommand(addPluginCmd)
}
