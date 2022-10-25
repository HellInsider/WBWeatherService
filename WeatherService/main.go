package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		defer f.Close()
		log.SetOutput(f)
	}

	closeCh := make(chan struct{})
	errChan := make(chan error)
	defer close(closeCh)
	defer close(errChan)
	go GErrGoroutinesHandler(closeCh, errChan)              //handles goroutines errors
	go AsyncWeatherUpdate(time.Second*60, closeCh, errChan) //updates weather info every args[0] time

	WriteCitiesToDB(errChan)  //Initial DB  cities filling (list of cities can be found in ./BashScripts/Cities.txt)
	WriteWeatherToDB(errChan) //Initial DB weather filling

	ServiceRun()
}
