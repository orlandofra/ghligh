/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"sync"

	"encoding/json"

	"github.com/orlandofra/ghligh/document"
	"github.com/orlandofra/ghligh/go-poppler"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "display info about pdf documents [json]",
	Long: `
	ghligh info file1.pdf file2.pdf [-i]

	shows information about pdf (author, publisher, modification date, etc...)
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

		infoChan := make(chan poppler.DocumentInfo)
		var wg sync.WaitGroup
		wg.Add(len(args))

		for _, arg := range args {
			go func(arg string) {
				defer wg.Done()
				doc, err := document.Open(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error opening %s: %v\n", arg, err)
					return
				}
				defer doc.Close()
				infoChan <- doc.Info()
			}(arg)
		}

		go func() {
			wg.Wait()
			close(infoChan)
		}()

		var infos []poppler.DocumentInfo
		for info := range infoChan {
			infos = append(infos, info)
		}

		var jsonBytes []byte
		if indent {
			jsonBytes, err = json.MarshalIndent(infos, "", "	")
		} else {
			jsonBytes, err = json.Marshal(infos)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	infoCmd.Flags().BoolP("indent", "i", false, "indent the json data")
}
