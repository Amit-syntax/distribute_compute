package main

import (
	"github.com/spf13/cobra"
	"fmt"
)

var rootCmd = &cobra.Command{
	Use:   "dc",
	Short: "A CLI tool for communicating with distribute_compute server",
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().StringP("server-ip", "ip", "", "distribute compute server IP")
	// Mark server-ip flag as required
	if err := connectCmd.MarkFlagRequired("server-ip"); err != nil {
		fmt.Println("Error marking server-ip as required:", err)
	}

	connectCmd.Flags().StringP("port", "p", "", "port of distribute_compute server")
	// Mark server-ip flag as required
	if err := connectCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println("Error marking port as required:", err)
	}

}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to a 'distribute compute' server",
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")
		fmt.Printf("Hello, %s!\n", name)

	},
}

func main() {
	rootCmd.Execute()
}

