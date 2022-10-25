package main

import (
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Полезные команды для работы с приложением.",
}

func init() {
	rootCmd.AddCommand(toolCmd)
}
