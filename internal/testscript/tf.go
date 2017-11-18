package main

import "github.com/gogits/logs"

func main() {
	log := logs.NewLogger(0)
	log.SetLogger("file", `{"filename":"test.log"}`)
	log.Info("who")
	log.Info("who")
	log.Info("who")
	log.Info("who")
	log.Critical("shit %d", 123)
	log.Info("who")
	log.Fatal("fatal")
}
