package main

import (
	"github.com/spf13/cobra"
	"fmt"

	"github.com/Amit-syntax/distribute_compute/internal/client"
)

var rootCmd = &cobra.Command{
	Use:   "dc",
	Short: "A CLI tool for communicating with distribute_compute server",
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringP("host", "", "", "distribute compute server host")
	// Mark host flag as required
	if err := connectCmd.MarkFlagRequired("host"); err != nil {
		fmt.Println("Error marking host as required:", err)
	}

	connectCmd.Flags().StringP("port", "p", "", "port of distribute_compute host")
	// Mark port flag as required
	if err := connectCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println("Error marking port as required:", err)
	}

}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to a 'distribute compute' server",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")

		client.Connect(host, port)

	},
}

func main() {
	rootCmd.Execute()
}

