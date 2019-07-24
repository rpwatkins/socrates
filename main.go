package main

import (
	"github.com/rpwatkins/socrates/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	cmd.Execute()
}
