/*
Copyright © 2025 Francesco Orlando orlando@stronzi.org
*/
package cmd

import (
	"os"

	"github.com/orlandofra/ghligh/cmd/tag"
	"github.com/orlandofra/ghligh/go-poppler"
	"github.com/spf13/cobra"
)

var warnings bool

var rootCmd = &cobra.Command{
	Use:   "ghligh",
	Short: "pdf highlights swiss knife",
	Long:  `ghligh can be used to manipulate pdf files in various ways.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		if !warnings {
			poppler.DisablePopplerWarnings()
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(tag.TagCmd)
	rootCmd.PersistentFlags().BoolVar(&warnings, "warnings", false, "show poppler warnings")
}
