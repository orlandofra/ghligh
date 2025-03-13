/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"

	"crypto/sha256"
)

var inputFiles []string

type importedAnnots struct {
	// FIXME just internatl map[string]struct where
	//	struct contains map[string]bool and document.AnnotsMap
	internal     map[string]document.AnnotsMap
	annotsHashes map[string]map[string]bool
	mutex        sync.Mutex
}

func (ia *importedAnnots) get(hash string) document.AnnotsMap {
	return ia.internal[hash]
}

func (ia *importedAnnots) init(hash string) {
	ia.mutex.Lock()
	defer ia.mutex.Unlock()
	if ia.internal[hash] == nil {
		ia.internal[hash] = make(document.AnnotsMap)
	}
	if ia.annotsHashes[hash] == nil {
		ia.annotsHashes[hash] = make(map[string]bool)
	}
}

func (ia *importedAnnots) check(docHash string, annotsHash string) bool {
	ia.mutex.Lock()
	defer ia.mutex.Unlock()
	ok := ia.annotsHashes[docHash][annotsHash]
	return ok
}

func (ia *importedAnnots) insert(hash string, am document.AnnotsMap) error {
	amHash, err := hashAnnotsMap(am)
	if err != nil {
		return err
	}

	present := ia.check(hash, amHash)
	if !present {
		ia.mutex.Lock()
		for key, value := range am {
			ia.internal[hash][key] = append(ia.internal[hash][key], value...)
			ia.annotsHashes[hash][amHash] = true
		}
		ia.mutex.Unlock()
	}
	return nil
}

func hashAnnotsMap(am document.AnnotsMap) (string, error) {
	jsonBytes, err := json.Marshal(am)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("%x", h), nil
}

func loadImportedAnnots(ia *importedAnnots, reader io.Reader) {
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read data: %v\n", err)
		return
	}

	var importedDocs []document.GhlighDoc

	err = json.Unmarshal(data, &importedDocs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	for _, importedDoc := range importedDocs {
		hash := importedDoc.HashBuffer
		ia.init(hash)
		ia.insert(hash, importedDoc.AnnotsBuffer)
	}
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import highlights from json file",
	Long: `
	ghligh import foo.pdf bar.pdf ... [--from fnord.json] [--from kadio.json] [-0] [--save=false]

	will import into foo.pdf bar.pdf etc... the highlights from file specified
	with the --from flag

	if -0 is set ghligh will read json from stdin

	--save=false will run without saving documents, it will just tells you how
	many annotations from the json files specified will be imported
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		stdin, err := cmd.Flags().GetBool("stdin")
		if err != nil {
			cmd.Help()
			return
		}

		save, err := cmd.Flags().GetBool("save")
		if err != nil {
			cmd.Help()
			return
		}

		if stdin == false && len(inputFiles) == 0 {
			fmt.Fprintf(os.Stderr, "nowhere to put output I am not doing anything\n")
			return
		}

		// Load Annot Maps
		ia := importedAnnots{
			internal:     make(map[string]document.AnnotsMap),
			annotsHashes: make(map[string]map[string]bool),
		}

		var wg sync.WaitGroup

		wg.Add(len(inputFiles))
		for _, file := range inputFiles {
			//wg.Add(1)
			go func(path string) {
				defer wg.Done()

				f, err := os.Open(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not open file %s: %v\n", path, err)
					return
				}
				defer f.Close()

				loadImportedAnnots(&ia, f)
			}(file)
		}

		wg.Wait()

		if stdin {
			loadImportedAnnots(&ia, os.Stdin)
		}

		// load from inputFiles
		for _, file := range args {
			doc, err := document.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error loading %s: %v", file, err)
				continue
			}

			hash := doc.HashDoc()

			num, err := doc.Import(ia.get(hash))
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not import highlights into %s: %v\n", file, err)
			} else {
				fmt.Fprintf(os.Stderr, "imported %d annots into %s\n", num, file)
				if save {
					doc.Save()
				}
			}
			doc.Close()

		}

	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().BoolP("stdin", "0", false, "read json from stdin")
	importCmd.Flags().BoolP("save", "", true, "save the file with new annotation importer")
	importCmd.Flags().StringArrayVarP(&inputFiles, "from", "f", []string{}, "files to import annots from")
}
