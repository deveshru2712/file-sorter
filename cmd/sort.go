package cmd

import (
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:     "sort",
	Aliases: []string{"organize", "divide"},
	Short:   "organize files into different folder",
	Long:    "This command creates different folder based on different files that are available.",
	Args:    cobra.ExactArgs(1),
	// should handle error as well
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Files sorted")
		SortFile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}
