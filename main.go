package main

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Path string
}

func parseFlags() *Config {
	pathFlag := flag.String("path", "", "Path to files")
	flag.Parse()

	conf := &Config{
		Path: *pathFlag,
	}

	return conf
}

func usage() {
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	conf := parseFlags()

	if conf.Path == "" {
		usage()
	}

	repo := NewRepository(conf.Path)
	samples, err := repo.Samples()
	if err != nil {
		log.Fatal("Could not read samples from repository:", err)
	}

	for _, s := range samples {
		log.Println(s.Title, "G:", s.Group)
	}
}
