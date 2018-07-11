// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/hr1sh1kesh/vloadgenerator/src"
	"github.com/spf13/cobra"
)

var rate int
var duration int

// datagenCmd represents the datagen command
var datagenCmd = &cobra.Command{
	Use:   "datagen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: generateData,
}

func init() {
	rootCmd.AddCommand(datagenCmd)
	datagenCmd.PersistentFlags().IntVarP(&rate, "rate", "n", 5, "Request Rate")
	datagenCmd.MarkFlagRequired("rate")
	datagenCmd.PersistentFlags().IntVarP(&duration, "duration", "d", 50, "Duration")
	datagenCmd.MarkFlagRequired("duration")
}

func generateData(cmd *cobra.Command, args []string) {
	src.GenerateLoadData(rate, duration, api)
}
