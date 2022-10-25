package main

import (
	log "github.com/sirupsen/logrus"
)

func GErrGoroutinesHandler(cls chan struct{}, err chan error) {
	select {
	case _, ok := <-cls:
		if !ok {
			log.Info("Shutting down error handler")
			return
		}
	case e, ok := <-err:
		if !ok {
			log.Errorf("Error with errorChan occured")
			return
		} else {
			log.Errorf("Goroutine recived error: %s", e)
		}
	}
}
