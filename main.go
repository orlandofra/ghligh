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
	"fmt"
	"os"
)

func usage(){
	fmt.Fprintf(os.Stderr, "ghligh: the pdf highlight swiss knife \n")
	fmt.Fprintf(os.Stderr, "\nusage: %s <command> [args...]\n", os.Args[0])

	fmt.Fprintf(os.Stderr, "\nCommands:")
	fmt.Fprintf(os.Stderr, "\n  export:\n")
	fmt.Fprintf(os.Stderr, "    Export highlights from one or more PDF files to a JSON file.\n")
	fmt.Fprintf(os.Stderr, "    If no JSON file is specified, outputs JSON data to stdout.\n")
	fmt.Fprintf(os.Stderr, "    Each JSON will contain highlights of every pdf file inserted.\n")
	fmt.Fprintf(os.Stderr, "    Options:\n")
	fmt.Fprintf(os.Stderr, "      -f, --from    Specify the PDF file(s) to export highlights from.\n")
	fmt.Fprintf(os.Stderr, "      -t, --to      Specify the JSON file(s) to save highlights to.\n")
	fmt.Fprintf(os.Stderr, "    example usage: ``ghligh export -f file1.pdf -f file2.pdf # outputs to \n")
	fmt.Fprintf(os.Stderr, "                            stdout``\n")

	fmt.Fprintf(os.Stderr, "\n  import:\n")
	fmt.Fprintf(os.Stderr, "    Cluster highlights from one or more JSON files, \n")
	fmt.Fprintf(os.Stderr, "    check if highlights are missing inside the corresponding PDF files and \n")
	fmt.Fprintf(os.Stderr, "    import missing highlights into each PDF file.\n")
	fmt.Fprintf(os.Stderr, "    Options:\n")
	fmt.Fprintf(os.Stderr, "      -f, --from    Specify the JSON file(s) containing annotations.\n")
	fmt.Fprintf(os.Stderr, "      -t, --to	    Specify the PDF file(s) for importing highlights.\n")
	fmt.Fprintf(os.Stderr, "    example usage: ``ghligh import -f one_highlight_collection.json \n")
	fmt.Fprintf(os.Stderr, "                            -f another_highlight_collection.json -t file1.pdf\n")
	fmt.Fprintf(os.Stderr, "                            -t file2.pdf -t file3.pdf``\n")


	fmt.Fprintf(os.Stderr, "\n  cat:\n")
	fmt.Fprintf(os.Stderr, "    Print the highlights from the given PDF file(s) to stdout.\n")
	fmt.Fprintf(os.Stderr, "    example usage: ``ghligh cat file1.pdf file2.pdf ...`` \n")

	fmt.Fprintf(os.Stderr, "\n  ls:\n")
	fmt.Fprintf(os.Stderr, "    Show pdf files that contains at least one highlight, for directory\n")
	fmt.Fprintf(os.Stderr, "    show files inside directory that contains highlights.\n")


	fmt.Fprintf(os.Stderr, "    example usage: ``ghligh ls file1.pdf directory file2.pdf ...``\n")
	fmt.Fprintf(os.Stderr, "    example usage: ``ghligh ls # show pdf files with highlights in current\n")
	fmt.Fprintf(os.Stderr, "    working directory``\n\n\n")

	os.Exit(1)
}



func main(){
	if len(os.Args) < 2 {
		usage()
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "export":
		commandExport(args)
	case "import":
		commandImport(args)
	case "cat":
		commandCat(args)
	case "ls":
		commandLs(args)
	default:
		usage()
	}

}
