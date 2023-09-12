package main

import (
	"go-canal/libs/log"
	"go-canal/libs/service"
	"os"
	ossignal "os/signal"
	"syscall"
)

var (
	configFile string
	isColor    bool
)

func init() {
	//init.go
}

func main() {
	defer os.Exit(0)

	service.Run()
	defer service.Stop()

	signal := make(chan os.Signal, 1)
	ossignal.Notify(signal, os.Kill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Warnf("Application catch os signal: %s \n", (<-signal).String())
}
