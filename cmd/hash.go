/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"sync"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "display the ghligh hash used to identify a documet [json]",
	Long: `the ghligh hash is used to identify documents with different filenames / annotations and it is calculated using the text of some pages.

		ghligh hash file1.json file2.json [-i]

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

		//		var hashes []document.GhlighDoc
		hashChan := make(chan document.GhlighDoc)
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

				// A little hacky, set hash after closing the document
				//doc.Close()
				doc.HashBuffer = doc.HashDoc()

				hashChan <- *doc
			}(arg)
		}

		go func() {
			wg.Wait()
			close(hashChan)
		}()

		var hashes []document.GhlighDoc
		for doc := range hashChan {
			hashes = append(hashes, doc)
			doc.Close()
		}

		var jsonBytes []byte
		if indent {
			jsonBytes, err = json.MarshalIndent(hashes, "", "	")
		} else {
			jsonBytes, err = json.Marshal(hashes)
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))

	},
}

func init() {
	rootCmd.AddCommand(hashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	hashCmd.Flags().BoolP("indent", "i", false, "indent the json data")
}
