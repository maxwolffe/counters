/*
Copyright © 2021 Max Wolffe <max.alan.wolffe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const CountersKey = "counters"
const TimeLayout = time.UnixDate

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "counters",
	Short: "A CLI for managing timers",
	Long: `CLI for managing timers. For example:

$ counters create -n eggs
Creating a new timer with name: eggs

$ counters list
Listing counters...
- Counter: eggs, Duration: 34m39.788915s
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.counters.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFilePath := cfgFile
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cfgFilePath = home
		cobra.CheckErr(err)

		// Search config in home directory with name ".counters" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".counters")
	}

	viper.SetDefault(CountersKey, make(map[string]string))

	viper.AutomaticEnv() // read in environment variables that match

	var configNotFoundError viper.ConfigFileNotFoundError
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// TODO move this to debug logging
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if errors.As(err, &configNotFoundError) {
		fmt.Printf("Creating empty config at default location: %s", cfgFilePath)
		viper.SafeWriteConfig()
	}
}
