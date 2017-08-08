// +build darwin linux

package config

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
)

// LogWriter will return an io.writer that is created based on environment
// configuration
func LogWriter() io.Writer {
	var writers []io.Writer

	for _, log := range Cfg.LogLocations {
		switch log {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "syslog":
			sl, err := syslog.New(6, "PRAELATUS")
			if err != nil {
				fmt.Println(err)
				continue
			}

			writers = append(writers, sl)
		default:
			var f *os.File
			var e error

			if _, err := os.Stat(log); err == nil {
				f, e = os.Open(log)
			} else {
				f, e = os.Create(log)
			}

			if e != nil {
				fmt.Printf("Error opening %s: %s", log, e.Error())
				continue
			}

			writers = append(writers, f)
		}
	}

	return io.MultiWriter(writers...)
}
