/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

var recursive bool

type resolver struct {
	paths   []string
	recurse bool
	ctx     context.Context
	ch      chan<- string
	wg      sync.WaitGroup
}

func (r *resolver) resolve() {
	for _, path := range r.paths {
		r.wg.Add(1)
		go r.resolvePath(path)
	}

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()
}

func (r *resolver) resolvePath(path string) {
	defer r.wg.Done()
	if err := r.ctx.Err(); err != nil {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		str, err := filepath.Abs(path)
		if err != nil {
			return
		}

		r.ch <- str
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving info for %s: %v\n", entry.Name(), err)
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		if info.IsDir() {
			if r.recurse {
				if err := r.ctx.Err(); err != nil {
					return
				}
				r.wg.Add(1)
				go r.resolvePath(fullPath)
			}
		} else if info.Mode().IsRegular() {
			r.ch <- fullPath
		}
	}
}

// ArgsOrCWD returns the provided args slice if non-empty.
// If args is empty, it returns a slice containing the current working directory.
//
// Example:
//
//	ArgsOrCWD([]string{"path1", "file2"}) -> []string{"path1", "file2"}
//	ArgsOrCWD([]string{}) -> []string{"."}
func ArgsOrCWD(args []string) []string {
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
			return nil
		}
		return []string{cwd}

	}
	return args
}

// returns true if file specified in path start with the PDF magic header
func isPDF(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		// We only care about printing permissions errors
		// and things like that
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return false
	}

	defer file.Close()

	header := make([]byte, 5)
	_, err = file.Read(header)
	if err != nil {
		return false
	}

	return bytes.Equal(header, []byte("%PDF-"))
}

// returns true if file contains at least one highlight or is tagged with "ls" (is ls-able by ghligh)
func HasHighlights(path string) bool {
	// Ensure is a pdf file, might block for other kind of files
	if !isPDF(path) {
		return false
	}

	doc, err := document.Open(path)
	if err != nil {
		return false
	}
	defer doc.Close()

	return doc.HasHighlights()
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "show files with highlights or tagged with 'ls' [unix]",
	Long: `
	ghligh ls file1.pdf directory  [-R] [-c]

	will show every file inside directory that contains highlights or it is marked
	ls with the ghligh tag add command

	ghligh ls # show files in current dir
	ghligh ls file1.pdf # if it outpus file1.pdf it means that file1.pdf contains highlights
	ghligh ls -c file1.pdf # same as ghligh ls file1.pdf but exit status will not be if file1.pdf doesnt
				contains highlights

	ghligh ls -R # do it recursively, be careful with symlink dir cycles, as I am to lazy to
			address that particular issue
`,
	Run: func(cmd *cobra.Command, args []string) {
		files := ArgsOrCWD(args)
		ch := make(chan string)
		ctx := context.Background()

		res := resolver{
			paths:   files,
			recurse: recursive,
			ctx:     ctx,
			ch:      ch,
		}

		go res.resolve()

		var wg sync.WaitGroup
		var found bool

		for file := range ch {
			wg.Add(1)
			go func(f string) {
				defer wg.Done()
				if HasHighlights(f) {
					found = true
					fmt.Printf("%s\n", f)
				}
			}(file)
		}

		wg.Wait()

		check, err := cmd.Flags().GetBool("check")
		if err != nil {
			cmd.Help()
			return
		}
		if check && !found {
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolVarP(&recursive, "recursive", "R", false, "List recursively")
	lsCmd.Flags().BoolP("check", "c", false, "exit status is 1 if no file its found")
	// order pdf by time of something (modification / creation) ???
	//lsCmd.Flags().BoolP("time", "t", false, "ls by time")
}
