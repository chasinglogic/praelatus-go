package cli

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/praelatus/praelatus/api"
	"github.com/praelatus/praelatus/api/middleware"
	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/models"
	"github.com/praelatus/praelatus/store"
	"github.com/tylerb/graceful"
	"github.com/urfave/cli"
)

// this is only used when running in dev mode to make testing the api easier
func disableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Authorization,Token")
			w.Header().Add("Access-Control-Expose-Headers", "*")
			w.Header().Add("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.Write([]byte{})
				return
			}

			next.ServeHTTP(w, r)
		})
}

func alwaysAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			u, _ := models.NewUser("testadmin", "test",
				"Test Testerson", "test@example.com", true)
			_ = middleware.SetUserSession(*u, w)
			r.Header.Set("Authorization", w.Header().Get("Token"))
			next.ServeHTTP(w, r)
		})
}

func runServer(c *cli.Context) error {
	log.SetOutput(config.LogWriter())

	log.Println("Starting Praelatus...")
	log.Println("Initializing database...")

	s := config.Store()
	ss := config.SessionStore()

	if sql, ok := s.(store.Migrater); ok {
		log.Println("Migrating database...")
		err := sql.Migrate()
		if err != nil {
			log.Println("Error migrating database:", err)
			os.Exit(1)
		}
	}

	log.Println("Prepping API")
	r := api.New(s, ss)
	if c.Bool("devmode") || os.Getenv("PRAELATUS_DEV_MODE") == "1" {
		log.Println("Running in dev mode, disabling cors and authentication...")
		r = disableCors(r)
		r = alwaysAuth(r)
	}

	log.Println("Listening on", config.Port())
	return graceful.RunWithErr(config.Port(), time.Minute, r)
}
