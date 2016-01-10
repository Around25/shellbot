package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var AppConfig *viper.Viper

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "shellbot",
	Short: "Dead simple provisioning tool",
	Long:  `Dead simple provisioning tool`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shellbot.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	AppConfig = viper.New()
	fmt.Printf("%v\n", cfgFile)
	if cfgFile != "" { // enable ability to specify config file via flag
		AppConfig.SetConfigFile(cfgFile)
	}

	AppConfig.SetConfigName(".shellbot") // name of config file (without extension)
	AppConfig.AddConfigPath(".")         // adding the current working directory first search path
	AppConfig.AddConfigPath("$HOME")     // adding home directory as second search path
	AppConfig.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := AppConfig.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", AppConfig.ConfigFileUsed())
	} else {
		fmt.Printf("%v\n", err)
	}

}
