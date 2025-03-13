/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat prints highlights of pdf files [unix][json]",
	Long: `
	ghligh cat file1.pdf file2.pdf ... [--json] [-i]

	will show every highlights inside pdf files specified
	if --json is set the output will be in json format

	if -i is set the json output will be indented
`,
	Run: func(cmd *cobra.Command, args []string) {

		useJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			cmd.Help()
			return
		}

		indent, err := cmd.Flags().GetBool("indent")
		if err != nil {
			cmd.Help()
			return
		}

		jsonCat := make(map[string][]document.HighlightedText)

		// for every arg
		for _, arg := range args {
			doc, err := document.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			highlights := doc.Cat()
			if !useJSON {
				for _, highlight := range highlights {
					if highlight.Contents != "" {
						fmt.Printf("%s {{{%s}}}", highlight.Text, highlight.Contents)
					} else {
						fmt.Printf("%s", highlight.Text)
					}
				}
			} else {
				jsonCat[doc.Path] = highlights
			}

			doc.Close()
		}

		var jsonBytes []byte
		if indent {
			jsonBytes, err = json.MarshalIndent(jsonCat, "", "	")
		} else {
			jsonBytes, err = json.Marshal(jsonCat)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))
	},
}

func init() {
	rootCmd.AddCommand(catCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	catCmd.Flags().BoolP("json", "j", false, "print highlights as json")
	catCmd.Flags().BoolP("indent", "i", false, "print highlights as json")
}
