package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aliykh/reddit-feed/api/docs"
	"github.com/aliykh/reddit-feed/internal/bootstrap"
	"github.com/aliykh/reddit-feed/internal/config"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

// @title Reddit Feed Api
// @version 1.0
// @description REST API for reddit feed posts
// @contact.name Alloy
// @contact.email aliykhoshimov@gmail.com
// @license.name Toptal
// @license.url https://toptal.com
// @BasePath /api/v1
func main() {

	flag.Parse()

	// load application configurations
	cfg, err := config.Load(*flagConfig, nil)
	if err != nil {
		log.Fatalf("failed to load application configuration: %s\n", err)
	}

	//// Swagger config
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", cfg.ServerPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	interruptChan := make(chan os.Signal, 1)

	signal.Notify(interruptChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		oscall := <-interruptChan

		log.Printf("system call:%+v\n", oscall)

		close(interruptChan)
		signal.Stop(interruptChan)

		cancel()
	}()

	a := bootstrap.New(cfg, ctx)

	a.Run(ctx)

	os.Exit(0)

}
