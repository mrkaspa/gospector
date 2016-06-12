package main

import (
	"errors"
	"os"

	"encoding/json"

	"fmt"

	"io/ioutil"

	"github.com/codegangsta/cli"
)

const VERSION = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "gospector"
	app.Usage = "Check the README.md here httpds://github.com/mrkaspa/gospector"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Usage: "Directory to gospect",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "Config file for gospector",
		},
	}
	app.Action = func(c *cli.Context) error {
		err := run(c)
		if err != nil {
			fmt.Println(err)
		}
		return err
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	dir := c.String("dir")
	if dir == "" {
		dir, _ = os.Getwd()
	} else if info, err := os.Stat(dir); err != nil {
		return err
	} else if !info.IsDir() {
		return errors.New("The dir must be a directory valid")
	}

	configFile := c.String("config")
	if configFile == "" {
		configFile = dir + "/gospector.json"
	}
	if _, err := os.Stat(configFile); err != nil {
		return err
	}

	config, err := readConfig(configFile)
	if err != nil {
		return err
	}

	g := createGospector(dir, config)
	errors := g.execute()
	if len(errors) > 0 {
		fmt.Println("\n****WORDS FOUND****\n")
		for _, err := range errors {
			fmt.Println(err)
			fmt.Println()
		}
	} else {
		fmt.Println("\n<< O K >>\n")
	}
	return nil
}

func readConfig(configFile string) (*gospectorConf, error) {
	var config gospectorConf

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
