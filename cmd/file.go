/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
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
		indent, err := cmd.Flags().GetInt("indent")
		if err != nil {
			fmt.Println(err)
		}

		lang, err := cmd.Flags().GetString("language")
		if !parser.IsValidLang(lang) {
			log.Fatalf(
				"%s is not a valid language. Valid languages: %v",
				lang,
				parser.GetValidLangs(),
			)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			reader := bufio.NewReader(f)
			if err != nil {
				fmt.Println(path)
			}

			res := parser.ParseTypes(reader, parser.Lang(lang), indent, name)
			println(res)
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.Flags().IntP("indent", "i", 4, "Output indentation")
	fileCmd.Flags().
		StringP("language", "l", "go", fmt.Sprintf("Output to especified language (%s)", parser.GetValidLangs()))
	fileCmd.Flags().StringP("name", "n", "Main", "Struct/Interface name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
