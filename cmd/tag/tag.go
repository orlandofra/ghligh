/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package tag

import (
	"regexp"
	"github.com/spf13/cobra"
)


var tags []string

var regex, exact string

func formatRegex(r string, boundaries string) string {
	//return `\b` + r + `\b`
	return boundaries + r + boundaries
}

func regexSlice(regex string, slice []string) []string {
	if regex == "" {
		return slice
	}

	var newSlice []string
	re, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	for _, s := range(slice){
		if re.MatchString(s){
			newSlice = append(newSlice, s)
		}
	}

	return newSlice
}

// tagCmd represents the tag command
var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "manage pdf tags",
	Long: `a tag is a string you can attach to a pdf
.`,
}

func init() {
	//rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO add to "add" command
	// tagCmd.Flags().StringArrayVarP(&tags, "tag", "t", []string{}, "Tag da associare ai file (può essere usato più volte)")
}
