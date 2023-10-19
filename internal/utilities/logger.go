package utilities

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"sync"
	"time"
)

type AppLogger struct {
	LogPath     string
	LastLogDate string
	CompressLog bool
	DailyRotate bool
	Logger      *lumberjack.Logger
}

var Log *log.Logger

func (p *AppLogger) SetAppLogger() error {
	if p.LogPath == "" {
		return fmt.Errorf("unspecify log file path")
	}

	file, err := os.OpenFile(p.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Log = log.New(file, "", log.Ldate|log.Ltime)
	p.Logger = &lumberjack.Logger{
		Filename:  p.LogPath,
		Compress:  p.CompressLog,
		LocalTime: true,
	}

	Log.SetOutput(p.Logger)

	p.LastLogDate = time.Now().Format("2006-01-02")

	if p.DailyRotate == true {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			p.LogRotator()
		}()

	}

	return nil
}

// LogRotator :
func (p *AppLogger) LogRotator() {
	for {
		// If Logger not set then skip
		if p.Logger == nil {
			continue
		}
		// Rotate Log If LastLog Date != Current Date
		if p.LastLogDate != time.Now().Format("2006-01-02") {
			// Set LastLogDate to Current Date
			p.LastLogDate = time.Now().Format("2006-01-02")
			err := p.Logger.Rotate()
			if err != nil {
				log.Println("Daily Log Rotator Err: ", err.Error())
				return
			}
			log.Println("| Log Rotated")
		}
		// Sleep every 5 seconds
		time.Sleep(time.Duration(5) * time.Second)
	}
}
