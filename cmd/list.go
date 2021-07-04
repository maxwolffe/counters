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

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running timers",
	Long:  `List all currently running timers`,
	Run: func(cmd *cobra.Command, args []string) {
		countersMap := viper.GetStringMapString(CountersKey)
		if len(countersMap) == 0 {
			fmt.Println("No counters to list...")
		}
		for k, v := range countersMap {
			listedTime, err := time.Parse(TimeLayout, v)
			cobra.CheckErr(err)

			timeDifference := time.Now().Sub(listedTime)
			fmt.Printf("%s - %s\n", k, timeDifference)
		}
	},
}
