package main

import (
	"flag"
	"log"
	"os"

	"github.com/cameliot/alpaca"
	"github.com/cameliot/alpaca/meta"
)

var version = "unknown"

func usage() {
	flag.PrintDefaults()
	os.Exit(-1)
}

func main() {
	log.Println("Starting klang3 v.", version)

	// Initialize configuration
	config := parseFlags()

	if config.RepoPath == "" {
		usage()
	}

	// Initialize MQTT connection
	actions, dispatch := alpaca.DialMqtt(
		config.Mqtt.BrokerUri(),
		alpaca.Routes{
			"sampler": config.Mqtt.BaseTopic,
			"meta":    "v1/_meta",
		},
	)

	// Initialize repository
	repo := NewRepository(config.RepoPath)
	err := repo.Update()
	if err != nil {
		log.Fatal("Could not read samples from repository:", err)
	}

	samplerActions := make(alpaca.Actions)
	metaActions := make(alpaca.Actions)

	// Initialize Soundboard Service
	samplerSvc := NewSamplerSvc(repo)
	go samplerSvc.Handle(samplerActions, dispatch)

	// Hanlde meta actions for service discovery
	metaSvc := meta.NewMetaSvc(
		"sampler@mainhall",
		"klang3",
		version,
		"Klang3 MQTT SoundBoard",
	)
	go metaSvc.Handle(metaActions, dispatch)

	for action := range actions {
		samplerActions <- action
		metaActions <- action
	}

}
