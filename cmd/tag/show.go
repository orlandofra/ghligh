/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package tag

import (
	"fmt"
	"os"

	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/orlandofra/ghligh/document"
)

var tagShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show ghligh tags of pdf files [json]",
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

		indent, err := cmd.Flags().GetBool("indent")
		if err != nil {
			cmd.Help()
			return
		}

		exportTags := make(map[string][]string)
		for _, file := range(args){
			doc, err := document.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			if regex != "" {
				regex = formatRegex(regex, "")
			}
			tags := regexSlice(regex, doc.GetTags())

			exportTags[doc.Path] = tags

			doc.Close()

		}

		var jsonBytes []byte
		if indent {
			jsonBytes, err = json.MarshalIndent(exportTags, "", "	")
		} else {
			jsonBytes, err = json.Marshal(exportTags)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))
	},
}

func init() {
	TagCmd.AddCommand(tagShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showtagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	tagShowCmd.Flags().BoolP("indent", "i", false, "indent the json data")
	tagShowCmd.Flags().StringVarP(&regex, "regex", "r", "", "regex")
}
