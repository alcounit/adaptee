package adaptee

import (
	log "github.com/sirupsen/logrus"
)

//Configuration ....
type Configuration struct {
	SelenosisURL string
}

//App ...
type App struct {
	logger       *log.Logger
	selenosisURL string
}

//New ...
func New(logger *log.Logger, conf Configuration) *App {
	return &App{
		logger:       logger,
		selenosisURL: conf.SelenosisURL,
	}
}
