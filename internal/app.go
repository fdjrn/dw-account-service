package internal

import (
	"fmt"
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ExitGracefully() {
	// close mongodb connection
	if err := db.Mongo.Disconnect(); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("[CONN] db >> connection successfully disconnected")

	// close kafka connection
	_ = kafka.Producer.Close()
	log.Println("[CONN] kafka >> producer successfully closed")
}

// SetupCloseHandler :
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal,... Good Bye...")
		ExitGracefully()
		os.Exit(0)
	}()
}
