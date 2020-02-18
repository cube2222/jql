package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/nwidger/jsoncolor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cube2222/jql/jql/app"
)

var (
	cfgFile    string
	monochrome bool
)

type encoder interface {
	Encode(v interface{}) error
	SetIndent(prefix, indent string)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jql <query>",
	Short: "JSON Query Processor with a Lispy syntax.",
	Long:  string(MustAsset("../README.md")),
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := json.NewDecoder(bufio.NewReaderSize(os.Stdin, 4096*16))
		w := bufio.NewWriterSize(os.Stdout, 4096*16)
		defer w.Flush()
		var output encoder
		if monochrome {
			output = json.NewEncoder(w)
		} else {
			output = jsoncolor.NewEncoder(w)
		}
		output.SetIndent("", "  ")

		if len(args) == 0 {
			args = append(args, "(id)")
		}
		app := app.NewApp(args[0], input, output)

		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jql.yaml)")
	rootCmd.PersistentFlags().BoolVar(&monochrome, "monochrome", false, "monochrome (don't colorize output)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jql" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jql")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
