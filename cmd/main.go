package main

import (
	"fmt"
	"github.com/dw-account-service/configs"
	"github.com/dw-account-service/internal"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/kafka"
	"github.com/dw-account-service/internal/routes"
	"log"
	"sync"
)

func main() {
	var err error
	internal.SetupCloseHandler()

	defer internal.ExitGracefully()

	// Service Initialization
	err = configs.Initialize()
	if err != nil {
		log.Fatalln(fmt.Sprintf("error on config initialization: %s", err.Error()))
	}

	if err = db.Mongo.Connect(); err != nil {
		log.Fatalln(fmt.Sprintf("error on mongodb connection: %s", err.Error()))
	}

	wg := &sync.WaitGroup{}

	// Start Messages Producer
	wg.Add(1)
	go func() {
		err = kafka.Initialize()
		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()

	// Start Rest API
	wg.Add(1)
	go func() {
		err = routes.Initialize()
		if err != nil {
			log.Fatalln(err)
		}

		wg.Done()
	}()

	wg.Wait()

}
