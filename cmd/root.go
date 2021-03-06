package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool
var cfgFile string
var accessToken string
var baseURL string
var userID string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devflow",
	Short: "utils for working with TP",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devflow.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose")
	rootCmd.PersistentFlags().String("baseurl", "", "base URL to TP. Ex: https://project.tpondemand.com")
	rootCmd.PersistentFlags().String("accesstoken", "", "Your TP access token")
	rootCmd.PersistentFlags().String("userid", "", "Your TP user ID")
	viper.BindPFlag("baseurl", rootCmd.PersistentFlags().Lookup("baseurl"))
	viper.BindPFlag("accesstoken", rootCmd.PersistentFlags().Lookup("accesstoken"))
	viper.BindPFlag("userid", rootCmd.PersistentFlags().Lookup("userid"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".devflow" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(home + "/.config")
		viper.SetConfigName(".devflow")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	if viper.IsSet("accesstoken") {
		accessToken = viper.GetString("accesstoken")
	} else {
		fmt.Println("accesstoken is not set")
		os.Exit(1)
	}
	if viper.IsSet("baseurl") {
		baseURL = viper.GetString("baseurl")
	} else {
		fmt.Println("baseurl is not set")
		os.Exit(1)
	}
	if viper.IsSet("userid") {
		userID = viper.GetString("userid")
	} else {
		fmt.Println("userid is not set")
		os.Exit(1)
	}
}
