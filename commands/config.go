// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"os"

	"github.com/praelatus/praelatus/config"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(show)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Commands for interacting with the config file.",
}

var show = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Cfg)
	},
}

var gen = &cobra.Command{
	Use:   "gen",
	Short: "Generate a config file based on environment variables and defaults.",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open("config.json")
		if err != nil && !os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}

		defer f.Close()

		fmt.Println(config.Cfg)

		f, _ = os.Create("config.json")
		_, err = f.Write([]byte(config.Cfg.String()))
		if err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}
