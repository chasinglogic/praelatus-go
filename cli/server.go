package cli

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	mgo "gopkg.in/mgo.v2"

	"github.com/praelatus/backend/api"
	"github.com/praelatus/backend/api/middleware"
	"github.com/praelatus/backend/config"
	"github.com/praelatus/backend/models"
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

func connectDB() *mgo.Session {
	conn, err := mgo.Dial(config.DBURL())
	if err != nil {
		panic(err)
	}

	return conn
}

func runServer(c *cli.Context) error {
	log.SetOutput(config.LogWriter())

	log.Println("Starting Praelatus...")
	log.Println("Connecting to database...")
	r := api.New(connectDB())
	if c.Bool("devmode") || os.Getenv("PRAELATUS_DEV_MODE") == "1" {
		log.Println("Running in dev mode, disabling cors and authentication...")
		r = disableCors(r)
		r = alwaysAuth(r)
	}

	if c.Bool("profile") {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	log.Println("Listening on", config.Port())
	return graceful.RunWithErr(config.Port(), time.Minute, r)
}
