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

	"path/filepath"
)


func lsCheckFile(arg string) bool {
	doc, err := ghligh.Open(arg)
	if err != nil {
		return false
	}
	defer doc.Close()

	if doc.HasHighlights() {
		return true
	}

	return false

}

func lsArgs(args []string) []string{
	newArgs := make([]string, len(args))
	copy(newArgs, args)

	if len(newArgs) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return nil
		}

		newArgs = append(newArgs, dir)
	}
	return newArgs
}

func filesInsideDirectory(dir string) []string {
	/* UGLY
	* this is an easy implementation for `ls`
	* where if directory is given inside a list
	* of strings it just remove the directory and put
	* inside the files (just the files)
	*
	* for a better implementation (that should be able
	* to recurse, `-R`) it should have a check for
	* cycles with symbolic linked dirs and things like
	* that. Since this is a lazy written software
	* i left it minimal like this...
	*/
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []string

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files
}

func lsProcessArgs(args []string) []string{
	newArgs := make([]string, len(args))

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghligh: %v\n", err)
			continue
		}

		if info.IsDir() {
			files := filesInsideDirectory(arg)
			newArgs = append(newArgs, files...)
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	return newArgs
}

func commandLs(args []string){
	newArgs := lsArgs(args)
	newArgs = lsProcessArgs(newArgs)

	/* something like exitStatus := C.EXIT_SUCESS
	* would be sufficient to be diagnosed with
	* paranoid schizofrenia, so I left 1 and 0
	* as values
	*/
	exitStatus := 1

	for _, arg := range(newArgs){
		doc, err := ghligh.Open(arg)
		if err != nil {
			continue
		}
		if doc.HasHighlights(){
			exitStatus = 0
			/* absPath is better for scripting if we are using ghligh
			* to pipe the output for other programs
			*/
			absPath, _ := filepath.Abs(arg)
			fmt.Printf("%s\n", absPath)
		}
		doc.Close()
	}

	os.Exit(exitStatus)
}
