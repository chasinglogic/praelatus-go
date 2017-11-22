// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package commands

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Version of Praelatus
	Version string
	// Commit hash this build was built against
	Commit string
)

func init() {
	Root.AddCommand(server)
	Root.AddCommand(db)
	Root.AddCommand(admin)
	Root.AddCommand(versionCmd)
}

// Execute runs the root command
func Execute() {
	Root.Execute()
}

// Root is the global CLI instance
var Root = &cobra.Command{
	Use:   "praelatus",
	Short: "A free and open source ticketing system.",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(cmd, args)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information.",
	Run: func(cmd *cobra.Command, args []string) {
		var version string

		if strings.HasPrefix(Version, "SNAPSHOT") {
			version = "DEV-" + Commit
		} else {
			version = Version
		}

		fmt.Printf("Praelatus %s#%s %s/%s\n",
			version, Commit, runtime.GOOS, runtime.GOARCH)
	},
}
