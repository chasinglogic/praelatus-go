// Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
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
)

// Config holds much of the configuration for praelatus, if reading from the
// configuration you should use the helper methods in this package as they do
// some prequisite processing and return appropriate types.
type Config struct {
	DBURL         string   `required:"true"`
	DBName        string   `required:"true"`
	DataDirectory string   `required:"true"`
	Port          string   `required:"true"`
	LogLocations  []string `required:"true"`
	InstanceName  string   `required:"true"`
}

// Public returns a safe for public consumption Config
func (c Config) Public() interface{} {
	return Config{}
}

func (c Config) String() string {
	b, e := json.MarshalIndent(c, "", "\t")
	if e != nil {
		return ""
	}

	return string(b)
}

// Save writes the config to a json file
func (c Config) Save() error {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("config.json", b, 0600)
}

// Cfg is the global config variable used in the helper methods of this package
var Cfg Config

func init() {
	Cfg.DBURL = os.Getenv("PRAELATUS_DB_URL")
	if Cfg.DBURL == "" {
		Cfg.DBURL = "mongodb://localhost:27017/praelatus"
	}

	Cfg.DataDirectory = os.Getenv("PRAELATUS_DATA_DIRECTORY")
	if Cfg.DataDirectory == "" {
		path, err := filepath.Abs("data")
		if err != nil {
			Cfg.DataDirectory = "data"
		} else {
			Cfg.DataDirectory = path
		}
	}

	Cfg.DBName = os.Getenv("PRAELATUS_DB")
	if Cfg.DBName == "" {
		Cfg.DBName = "praelatus"
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

	Cfg.InstanceName = os.Getenv("PRAELATUS_INSTANCE_NAME")
	if Cfg.InstanceName == "" {
		Cfg.InstanceName = "Praelatus"
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

// WebWorkers returns the number of web workers to run for sending http
// requests from hooks
func WebWorkers() int {
	return 10
}

var Logger *log.Logger

func DataDir() string {
	return Cfg.DataDirectory
}
