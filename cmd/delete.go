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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteName string

func init() {
	deleteCmd.Flags().StringVarP(&deleteName, "name", "n", "", "The name of the timer to delete.")
	deleteCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(deleteCmd)
}

// This does a pretty gross thing of mutating the input parameter.
func removeFromMap(countersMap map[string]string, keyToRemove string) {
	_, ok := countersMap[keyToRemove]
	if !ok {
		fmt.Printf("No timer with name %s in timer list", keyToRemove)
		os.Exit(1)
	}
	fmt.Printf("Deleting timer with name: %s", keyToRemove)
	delete(countersMap, keyToRemove)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a timer",
	Long:  `Deletes a given timer`,
	Run: func(cmd *cobra.Command, args []string) {
		countersMap := viper.GetStringMapString(CountersKey)
		removeFromMap(countersMap, deleteName)
		viper.Set(CountersKey, countersMap)
		err := viper.WriteConfig()
		cobra.CheckErr(err)
	},
}
