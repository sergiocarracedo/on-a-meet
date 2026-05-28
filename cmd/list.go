package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available camera devices",
	Long:  `Enumerates /dev/video* devices and shows driver information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list: not yet implemented (Phase 2)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
