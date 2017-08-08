package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/praelatus/backend/config"
)

func showConfig(c *cli.Context) error {
	fmt.Println(config.Cfg)
	return nil
}

func genConfig(c *cli.Context) error {
	f, err := os.Open("config.json")
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	fmt.Println(config.Cfg)

	if os.IsNotExist(err) {
		f, _ = os.Create("config.json")
		_, err = f.Write([]byte(config.Cfg.String()))
		return err
	}

	fmt.Println("config.json already exists exiting...")
	return nil
}
