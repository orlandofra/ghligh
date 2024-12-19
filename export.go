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

	"encoding/json"
)

func commandExport(args []string){
	fromFiles, toFiles := parse2Files(args)

	if fromFiles == nil {
		fmt.Fprintf(os.Stderr, "ghligh: please specify at least one pdf file via -f or --from for export\n")
		os.Exit(1)
	}

	if !areValidPDFs(fromFiles){
		os.Exit(1)
	}


	var exportedDocs []interface{}
	for _, file := range(fromFiles){
		doc, err := ghligh.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghligh: error opening %s: %v\n", file, err)
			continue
		}

		exportedData, err := doc.Export()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghligh: error exporting highlights of file %s: %v\n", file, err)
			continue
		}
		exportedDocs = append(exportedDocs, exportedData)

		doc.Close()
	}

	jsonData, err := json.MarshalIndent(exportedDocs, "", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ghligh: could not Indent %v\n", err)
		os.Exit(1)
	}

	if toFiles == nil {
		fmt.Printf(string(jsonData))
	} else {
		for _, file := range(toFiles){
			 saveToJSON(jsonData, file)
		}
		/* for file in toFiles, create file and save
		it as json */
	}

}

