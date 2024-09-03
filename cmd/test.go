package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "test",
	Short: "test print",
	Long:  `Test a new print`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello there")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
