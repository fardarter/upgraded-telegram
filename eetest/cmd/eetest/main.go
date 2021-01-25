package main

import "github.com/sirupsen/logrus"

var (
	buildTime string
	gitHash   string
)

var log = logrus.WithField("ctx", "main")

func main() {

	log.Info(buildTime)
	log.Info(gitHash)
}
