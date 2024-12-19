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
        "crypto/sha256"
        "io/ioutil"


        "os"
        "encoding/json"
)

type importedAnnotMaps struct {
	internal map[string]map[string]ghligh.AnnotsMap
}

func (iAM *importedAnnotMaps) Get(documentHash string) []ghligh.AnnotsMap {

	nestedMap, exists := iAM.internal[documentHash]
	if !exists {
		return nil
	}

	result := make([]ghligh.AnnotsMap, 0, len(nestedMap))
	for _, annotsMap := range nestedMap {
		result = append(result, annotsMap)
	}

	return result
}


func (iAM *importedAnnotMaps) loadFiles(files []string){
	for _, file := range(files){
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		var importedDocs []ghligh.GhlighDoc

		err = json.Unmarshal(data, &importedDocs)
		if err != nil {
			continue
		}
		for _, importedDoc := range importedDocs {
			documentHash := importedDoc.Hash
			annotsMapHash, err := iAM.hashAM(importedDoc.AnnotsBuffer)
			if err != nil {
				continue
			}

			if _, exists := iAM.internal[documentHash]; !exists {
				iAM.internal[documentHash] = make(map[string]ghligh.AnnotsMap)
			}
			iAM.internal[documentHash][annotsMapHash] = importedDoc.AnnotsBuffer
		}
	}
}

func NewImportedAnnotMaps() *importedAnnotMaps {
	return &importedAnnotMaps{
		internal: make(map[string]map[string]ghligh.AnnotsMap),
	}
}



func (iAM *importedAnnotMaps) hashAM(am ghligh.AnnotsMap) (string, error) {
	jsonBytes, err := json.Marshal(am)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("%x", h), nil
}


func commandImport(args []string){
	fromFiles, toFiles := parse2Files(args)
	exitStatus := 0


	if fromFiles == nil {
		fmt.Fprintf(os.Stderr, "ghligh: please specify at least one json file via -f or --from for import\n")
		os.Exit(1)
	}
	if toFiles == nil {
		fmt.Fprintf(os.Stderr, "ghligh: please specify at least one pdf file via -t or --to import\n")
		os.Exit(1)
	}

	if !areValidJSONFiles(fromFiles){
		os.Exit(1)
	}

	if !areValidPDFs(toFiles){
		os.Exit(1)
	}

	iam := NewImportedAnnotMaps()
	iam.loadFiles(fromFiles)

	for _, file := range(toFiles)	{
		doc, err := ghligh.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghligh: file %s is not a valid pdf\n", file)
			continue
		}
		hash := doc.HashDoc()

		ams := iam.Get(hash)
		for _, am := range(ams){
			n_imported, err := doc.Import(am)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ghligh: error %v\n", err)
				continue
			}
			if n_imported > 0 {
				fmt.Fprintf(os.Stderr, "ghligh: imported %d annotations in file %s\n", n_imported, file)
				if ok, err := doc.Save(); !ok{
					fmt.Fprintf(os.Stderr, "ghligh: could not save file: %v", err)
				}
			}
		}
		doc.Close()
	}

	os.Exit(exitStatus)

}
