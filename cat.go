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
)


func commandCat(args []string){
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "ghligh: please specify at least one file name for cat command\n")
		os.Exit(1)
	}

	exitStatus := 1

	for _, file := range(args){
		doc, err := ghligh.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ghligh: error opening %s: %v\n", file, err)
			continue
		}

		ch := make(chan string)

		go doc.Cat(ch)

		for annotText := range ch {
			exitStatus = 0
			fmt.Println(annotText)
		}

		doc.Close()
	}

	os.Exit(exitStatus)
}
