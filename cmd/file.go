/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fdanctl/jsontypify/internal/parser"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Args:  cobra.MinimumNArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		indent, err := cmd.Flags().GetInt("indent")
		if err != nil {
			fmt.Println(err)
		}
		lang, err := cmd.Flags().GetString("language")
		if !parser.IsValidLang(lang) {
			log.Fatalf("%s is not a valid language. Valid languages: %v", lang, parser.GetValidLangs())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(path, indent)
		res := parser.ParseTypes(data, parser.Lang(lang), indent)
		println(res)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.Flags().IntP("indent", "i", 4, "Output indentation")
	fileCmd.Flags().StringP("language", "l", "go", "Output to especified language (\"go\", \"ts\")")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
