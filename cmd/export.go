/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

var outputFiles []string

func writeJSONToFile(jsonBytes []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// IF FILE EXISTS ASK
	_, err = file.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export pdf highlights into json",
	Long: `
	ghligh export foo.pdf bar.pdf ... [--to fnord.json] [-1] [-i]

	will create one or more json file (specified with --to) or dump it
	to stdout (-1)

	-i will indent the json output
`,

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

		stdout, err := cmd.Flags().GetBool("stdout")
		if err != nil {
			cmd.Help()
			return
		}

		if !stdout && len(outputFiles) == 0 {
			fmt.Fprintf(os.Stderr, "nowhere to put output I am not doing anything\n")
			return
		}

		var exportedDocs []document.GhlighDoc
		for _, file := range args {
			doc, err := document.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error loading %s: %v", file, err)
				continue
			}

			doc.AnnotsBuffer = doc.GetAnnotsBuffer()
			doc.HashBuffer = doc.HashDoc()
			exportedDocs = append(exportedDocs, *doc)
		}

		var jsonBytes []byte
		if indent {
			jsonBytes, err = json.MarshalIndent(exportedDocs, "", "	")
		} else {
			jsonBytes, err = json.Marshal(exportedDocs)
		}
		if err != nil {
			panic(err)
		}

		for _, file := range outputFiles {
			err := writeJSONToFile(jsonBytes, file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}

		if stdout {
			fmt.Printf("%s\n", string(jsonBytes))
		}

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO flag toFiles
	exportCmd.Flags().BoolP("indent", "i", false, "indent the json data")
	exportCmd.Flags().BoolP("stdout", "1", false, "dump to stdout")

	exportCmd.Flags().StringArrayVarP(&outputFiles, "to", "t", []string{}, "files to save exported annots")
}
