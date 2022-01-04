package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lbhdc/block/pkg/v0/block"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	configPath = flag.String("config", "examples/config.toml", "path to config.toml")
	verbose    = flag.Bool("v", false, "verbose logging")
)

func main() {
	flag.Parse()
	config := block.NewConfigurationFromFile(*configPath)
	if err := config.Valid(); err != nil {
		log.WithError(err).Fatal("config.Valid")
	}
	if *verbose {
		fmt.Println(config)
	}
	server := block.NewServer(context.Background())
	server.Configure(config)
	defer func(server *block.Server) {
		if err := server.Shutdown(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}(server)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
