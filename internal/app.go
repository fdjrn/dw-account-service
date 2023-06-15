package internal

import (
	"github.com/dw-account-service/internal/db"
	"github.com/dw-account-service/internal/kafka"
	"github.com/dw-account-service/pkg/xlogger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ExitGracefully() {
	// close mongodb connection
	xlogger.Log.SetPrefix("[EXIT-APP] ")
	if err := db.Mongo.Disconnect(); err != nil {
		log.Println(err.Error())
		return
	}
	xlogger.Log.Println("| db connection successfully closed")

	// close kafka connection
	_ = kafka.Producer.Close()
	xlogger.Log.Println("| kafka producer successfully closed")
}

// SetupCloseHandler :
func SetupCloseHandler() {
	xlogger.Log.SetPrefix("[EXIT-APP] ")
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		xlogger.Log.Println("| Ctrl+C pressed in Terminal,... Good Bye...")
		ExitGracefully()
		os.Exit(0)
	}()
}
