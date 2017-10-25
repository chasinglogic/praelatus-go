// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package commands holds all the logic for the CLI of Praelatus
package commands

import (
	"log"
	"net/http"
	"os"
	"time"

	// Allows us to run profiling when flag is given
	_ "net/http/pprof"

	"github.com/praelatus/praelatus/api"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/spf13/cobra"
	"github.com/tylerb/graceful"
)

var (
	devMode     bool
	disableCORS bool
	profile     bool
)

func init() {
	server.Flags().BoolVar(&disableCORS, "nocors", false,
		"If given all Access-Control headers will be set to *")
	server.Flags().BoolVarP(&devMode, "dev-mode", "d", false,
		"Disables CORS and Authentication checks.")
	server.Flags().BoolVar(&profile, "profile", false,
		"Enables server performance profiling on localhost:6060")
}

var server = &cobra.Command{
	Use:   "serve",
	Short: "Run the praelatus API and UI server.",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(config.LogWriter())

		log.Println("Starting Praelatus...")
		log.Println("Connecting to database...")
		repo := config.LoadRepo()
		cache := config.LoadCache()

		api.Version = Version
		api.Commit = Commit

		r := api.New(repo, cache)
		if devMode || os.Getenv("PRAELATUS_DEV_MODE") != "" {
			log.Println("Running in dev mode, disabling cors and authentication...")
			r = disableCors(r)
			r = alwaysAuth(r)
		}

		if disableCORS {
			r = disableCors(r)
		}

		if profile {
			go func() {
				log.Println(http.ListenAndServe("localhost:6060", nil))
			}()
		}

		log.Println("Listening on", config.Port())
		err := graceful.RunWithErr(config.Port(), time.Minute, r)
		if err != nil {
			log.Println("Exited with error:", err)
		}
	},
}

// this is only used when running in dev mode to make testing the api easier
func disableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:4000")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Add("Access-Control-Expose-Headers", "X-Praelatus-Token, Content-Type, Authorization")
			w.Header().Add("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.Write([]byte{})
				return
			}

			next.ServeHTTP(w, r)
		})
}

func alwaysAuth(next http.Handler) http.Handler {
	u, _ := models.NewUser("testadmin", "test",
		"Test Testerson", "test@example.com", true)
	u.Roles = []models.UserRole{{Project: "TEST", Role: "Administrator"}}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_ = middleware.SetUserSession(*u, w)
			r.Header.Set("Authorization", w.Header().Get("Token"))
			next.ServeHTTP(w, r)
		})
}
