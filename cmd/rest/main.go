package main

import (
	"github.com/dw-account-service/internal"
	"github.com/dw-account-service/internal/conn/kafka"
	"github.com/dw-account-service/internal/routes"
	"log"
	"sync"
)

func main() {

	defer internal.ExitGracefully()

	var err error

	wg := &sync.WaitGroup{}

	// Config Initialization
	if err = internal.Initialize(); err != nil {
		log.Fatalln(err)
	}

	// Start Messages Producer
	wg.Add(1)
	go func() {
		err = kafka.StartMessageProducer()
		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()

	// Start Rest API
	wg.Add(1)
	go func() {
		err = routes.Start()
		if err != nil {
			log.Fatalln(err)
		}

		wg.Done()
	}()

	wg.Wait()

}
