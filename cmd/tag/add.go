/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package tag

import (
	"fmt"
	"io"
	"os"

	"github.com/orlandofra/ghligh/document"
	"github.com/spf13/cobra"
)

func readStdin() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// tagAddCmd represents the tag command
var tagAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a ghligh tag to a pdf file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		if stdin {
			tags = append(tags, readStdin())
		}
		if len(tags) == 0 {
			fmt.Fprintf(os.Stderr, "Either --tag or --stdin is required\n", err)
			os.Exit(1)
		}

		for _, file := range args {
			doc, err := document.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			for _, tag := range tags {
				doc.Tag(tag)
				fmt.Fprintf(os.Stderr, "added tag {{{%s}}} to %s\n", tag, file)
			}
			doc.Save()
			doc.Close()
		}
	},
}

func init() {
	TagCmd.AddCommand(tagAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO add to "add" command
	tagAddCmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Tag da associare ai file (può essere usato più volte)")
	tagAddCmd.Flags().BoolP("stdin", "0", false, "read tag from stdin")

	//if err := tagAddCmd.MarkFlagRequired("tag"); err != nil {
	//panic(err)
	//}

}
