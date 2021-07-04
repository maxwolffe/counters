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
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var name string

func init() {
	createCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the timer to create.")
	createCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new timer",
	Long:  `Create's a new timer starting immediately.`,
	Run: func(cmd *cobra.Command, args []string) {
		countersMap := viper.GetStringMapString(CountersKey)
		fmt.Printf("Creating a new timer with name: %s", name)
		countersMap[name] = time.Now().Format(TimeLayout)
		viper.Set(CountersKey, countersMap)
		err := viper.WriteConfig()
		cobra.CheckErr(err)
	},
}
