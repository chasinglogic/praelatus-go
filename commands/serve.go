// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package commands holds all the logic for the CLI of Praelatus
package commands

import (
	"log"
	"net/http"
	"time"

	// Allows us to run profiling when flag is given
	_ "net/http/pprof"

	"github.com/praelatus/praelatus/api"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/events"
	"github.com/spf13/cobra"
	"github.com/tylerb/graceful"
)

var (
	disableCORS bool
	profile     bool
)

func init() {
	server.Flags().BoolVar(&disableCORS, "nocors", false,
		"If given all Access-Control headers will be set to *")
	server.Flags().BoolVar(&profile, "profile", false,
		"Enables server performance profiling on localhost:6060")
}

var server = &cobra.Command{
	Use:   "serve",
	Short: "Run the praelatus API and UI server.",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(config.LogWriter())
		log.SetPrefix("[MAIN] ")

		log.Println("Starting Praelatus...")
		log.Println("Connecting to database...")
		rpo := loadRepo()

		api.Version = Version
		api.Commit = Commit

		r := api.New(rpo, middleware.Default)

		if profile {
			go func() {
				log.Println(http.ListenAndServe("localhost:6060", nil))
			}()
		}

		log.Println("Staring event manager...")
		go events.Run()

		log.Println("Listening on", config.Port())
		err := graceful.RunWithErr(config.Port(), time.Minute, r)
		if err != nil {
			log.Println("Exited with error:", err)
		}
	},
}
