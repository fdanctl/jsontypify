/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fdanctl/jsontypify/internal/parser"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jsontypify [flags] <file_path>",
	Args:  cobra.ExactArgs(1),
	Short: "Convert raw JSON to Go struct or TypeScript interface",
	Long:  `A tool to quickly convert raw JSON data into a Go struct or TypeScript interface.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var inputReader io.Reader

		
		if len(args) == 0 {
			cmd.Help()
			return
		} else if args[0] != "-" {
			// In Unix and CLI convention, a single dash "-" often represents stdin instead of a filename.
			file, err := os.Open(args[0])
			if err != nil {
				panic(err)
			}
			defer file.Close()
			inputReader = file
		} else {
			inputReader = cmd.InOrStdin()
		}

		indent, err := cmd.Flags().GetInt("indent")
		if err != nil {
			fmt.Println(err)
		}

		lang, err := cmd.Flags().GetString("language")
		if err != nil {
			fmt.Println(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}

		if !parser.IsValidLang(lang) {
			log.Fatalf(
				"%s is not a valid language. Valid languages: %s",
				lang,
				parser.GetValidLangs(),
			)
		}

		res := parser.ParseTypes(inputReader, parser.Lang(lang), indent, name)
		fmt.Fprintln(cmd.OutOrStdout(), res)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jsontypify.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().IntP("indent", "i", 4, "output indentation")

	langHelpMsg := fmt.Sprintf("output to especified language (%s)", parser.GetValidLangs())
	rootCmd.Flags().StringP("language", "l", "go", langHelpMsg)

	rootCmd.Flags().StringP("name", "n", "Main", "struct/interface name")
}
