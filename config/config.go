// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package config defines the config struct and provides utility methods for
// querying it. Additionally it handles loading the config.json if present and
// holds various other global values used in the app.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/praelatus/praelatus/repo"
	"github.com/praelatus/praelatus/repo/mongo"
)

type AWSConfig struct {
	Region  string
	BaseURL *string
}

// Config holds much of the configuration for praelatus, if reading from the
// configuration you should use the helper methods in this package as they do
// some prequisite processing and return appropriate types.
type Config struct {
	DBURL        string
	DBName       string
	SessionURL   string
	Port         string
	LogLocations []string
	SessionStore string
	AWS          AWSConfig
}

func (c Config) String() string {
	b, e := json.MarshalIndent(c, "", "\t")
	if e != nil {
		return ""
	}

	return string(b)
}

// Cfg is the global config variable used in the helper methods of this package
var Cfg Config

func init() {
	Cfg.DBURL = os.Getenv("PRAELATUS_DB_URL")
	if Cfg.DBURL == "" {
		Cfg.DBURL = "mongodb://localhost:27017/praelatus"
	}

	Cfg.DBName = os.Getenv("PRAELATUS_DB")
	if Cfg.DBName == "" {
		Cfg.DBName = "praelatus"
	}

	Cfg.SessionStore = os.Getenv("PRAELATUS_SESSION")
	if Cfg.SessionStore == "" {
		Cfg.SessionStore = "bolt"
	}

	Cfg.SessionURL = os.Getenv("PRAELATUS_SESSION_URL")
	if Cfg.SessionURL == "" {
		Cfg.SessionURL = "sessions.db"
	}

	Cfg.Port = os.Getenv("PRAELATUS_PORT")
	if Cfg.Port == "" {
		Cfg.Port = ":" + os.Getenv("PORT")
	}

	if Cfg.Port == ":" {
		Cfg.Port = ":8080"
	}

	Cfg.LogLocations = strings.Split(os.Getenv("PRAELATUS_LOGLOCATIONS"), ";")
	if os.Getenv("PRAELATUS_LOGLOCATIONS") == "" {
		Cfg.LogLocations = []string{"stdout"}
	}

	f, err := os.Open("config.json")
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	if os.IsNotExist(err) {
		return
	}

	jsn, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var c Config

	err = json.Unmarshal(jsn, &c)
	if err != nil {
		fmt.Println("Error unmarshaling config file:", err)
		os.Exit(1)
	}

	Cfg = c

	Logger = log.New(LogWriter(), "", log.LstdFlags)
}

// DBURL will return the environment variable PRAELATUS_DB if set, otherwise
// return the default development database url.
func DBURL() string {
	return Cfg.DBURL
}

// DBName will return the appropriate database name
func DBName() string {
	return Cfg.DBName
}

// Port will return the port / interfaces for the api to listen on based on the
// configuration
func Port() string {
	return Cfg.Port
}

// SessionURL will get the url to use for redis or file location for boltdb
func SessionURL() string {
	return Cfg.SessionURL
}

// WebWorkers returns the number of web workers to run for sending http
// requests from hooks
func WebWorkers() int {
	return 10
}

func LoadRepo() repo.Repo {
	return mongo.New(DBURL())
}

func LoadCache() repo.Cache {
	return mongo.NewCache(DBURL())
}

var Logger *log.Logger

func DataDir() string {
	path, err := filepath.Abs("data")
	if err != nil {
		return "data"
	}

	return path
}
