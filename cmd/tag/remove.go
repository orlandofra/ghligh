/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package tag

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/orlandofra/ghligh/document"
)

var tagRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove ghligh tags from a pdf files using regex",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		if regex == "" && exact == "" {
			fmt.Fprintf(os.Stderr, "either regex or exact must be set with --regex or --exact\n")
			os.Exit(1)
		}

		nosafe, err := cmd.Flags().GetBool("nosafe")
		if err != nil {
			cmd.Help()
			return
		}

		// just a little hack -> boundaries = nosafe ? "" : `\b`
		boundaries := map[bool]string{true: "", false: `\b`}[nosafe]
		regex = formatRegex(regex, boundaries)

		// if exact set overwrite regex
		if exact != "" {
			regex = `^` + exact + `$`
		}

		for _, file := range(args){
			doc, err := document.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			tags := regexSlice(regex, doc.GetTags())
			removedTags := doc.RemoveTags(tags)
			doc.Save()
			doc.Close()

			fmt.Printf("removed %d tags from %s\n", removedTags, doc.Path)
		}

	},
}

func init() {
	TagCmd.AddCommand(tagRemoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removetagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removetagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tagRemoveCmd.Flags().StringVarP(&regex, "regex", "r", "", "regex")
	tagRemoveCmd.Flags().StringVarP(&exact, "exact", "e", "", "exact")
	tagRemoveCmd.Flags().BoolP("nosafe", "", false, "don't use safe boundaries around regex")

	//if err := tagRemoveCmd.MarkFlagRequired("regex"); err != nil {
		//panic(err)
	//}
}
