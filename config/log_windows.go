package config

import (
	"fmt"
	"io"
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
