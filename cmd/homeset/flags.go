package main

import (
	"flag"
)

type appFlags struct {
	configPath string
}

func parseFlags() appFlags {
	defer flag.Parse()

	fs := appFlags{}

	flag.StringVar(&fs.configPath, "config", "config.yaml", "path to configuration yaml file")

	return fs
}
