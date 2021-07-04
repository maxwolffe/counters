package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

func init() {
	createCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the timer to create.")
	createCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new timer",
	Long:  `Create's a new timer starting immediately.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Creating a new timer with name: %s", name)
	},
}
