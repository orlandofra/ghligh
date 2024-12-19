// Package ghligh, a pdf highlights swiss knife
//
// Copyright (c) 2024 Francesco Orlando
//
// This file is part of a program licensed under the GNU General Public License, version 2.
// You should have received a copy of the GNU General Public License along with this program.
// If not, see the LICENSE.md file in the root directory of this repository or visit
// <https://www.gnu.org/licenses/old-licenses/gpl-2.0.html>.
package main

import (
        "github.com/orlandofra/ghligh/document"

	"fmt"

	"os"
	"flag"

	"io/ioutil"
	"encoding/json"
)

func isValidPDF(file string) bool {
	doc, err := ghligh.Open(file)
	defer doc.Close()
	return err == nil
}

func areValidPDFs(files []string) bool {
	for _, file := range(files){
		if !isValidPDF(file){
			fmt.Fprintf(os.Stderr, "ghligh: file %s is not a valid pdf\n", file)
			return false
		}
	}
	return true
}

func isValidJSONFile(file string) bool {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	var js json.RawMessage
	return json.Unmarshal(fileContent, &js) == nil
}

func areValidJSONFiles(files []string) bool {
	for _, file := range(files){
		if !isValidJSONFile(file){
			fmt.Fprintf(os.Stderr, "ghligh: file %s is not a valid json file\n", file)
			return false
		}
	}
	return true
}

func parse2Files(args []string) ([]string, []string){
	/* UGLY
	* should encapsulate flag inside a struct
	* and set things accordingly, this is just
	* a lazy workaround for import and export subcommands
	*/
	fs := flag.NewFlagSet("ghligh", flag.ExitOnError)

	var fromFiles []string
	var toFiles []string

	fs.Func("f", "List of files for -f", func(s string) error {
		fromFiles = append(fromFiles, s)
		return nil
	})

	fs.Func("from", "List of files for --from", func(s string) error {
		fromFiles = append(fromFiles, s)
		return nil
	})

	fs.Func("t", "List of files for -t", func(s string) error {
		toFiles = append(toFiles, s)
		return nil
	})

	fs.Func("to", "List of files for --to", func(s string) error {
		toFiles = append(toFiles, s)
		return nil
	})

	fs.Parse(args)

	return fromFiles, toFiles
}

func saveToJSON(jsonData []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
